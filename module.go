package CADDY_FILE_SERVER

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/h2non/bimg"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

func init() {
	caddy.RegisterModule(&Middleware{})
	httpcaddyfile.RegisterHandlerDirective("image_processor", parseCaddyfile)
	httpcaddyfile.RegisterDirectiveOrder("image_processor", "before", "respond")
}

// OnFail represents the possible values for the "on_fail" directive.
type OnFail string

const (
	// OnFailAbort returns a 500 Internal Server Error to the client.
	OnFailAbort OnFail = "abort"

	// OnFailBypass forces the response to return the initial (unprocessed) image.
	OnFailBypass OnFail = "bypass"
)

// Middleware allow user to do image processing on the fly using libvips
// With simple queries parameters you can resize, convert, crop your served images
type Middleware struct {
	logger   *zap.Logger
	OnFail   OnFail           `json:"on_fail,omitempty"`
	Security *SecurityOptions `json:"security,omitempty"`
}

func (*Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.image_processor",
		New: func() caddy.Module { return new(Middleware) },
	}
}

func (m *Middleware) Provision(ctx caddy.Context) error {
	m.logger = ctx.Logger()

	// Set default configuration
	m.OnFail = cmp.Or(m.OnFail, OnFailBypass)
	if m.Security != nil {
		if err := m.Security.Provision(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *Middleware) Validate() error {
	switch m.OnFail {
	case OnFailAbort, OnFailBypass:
		// Valid values
	default:
		return fmt.Errorf("invalid value for on_fail: '%s' (expected 'abort', or 'bypass')", m.OnFail)
	}

	if m.Security != nil {
		if err := m.Security.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	//Automatic return if not options set
	if r.URL.RawQuery == "" {
		return next.ServeHTTP(w, r)
	}

	responseRecorder := caddyhttp.NewResponseRecorder(w, &bytes.Buffer{}, func(status int, header http.Header) bool {
		return true
	})

	if err := next.ServeHTTP(responseRecorder, r); err != nil {
		return err
	}

	if responseRecorder.Status() != 200 || responseRecorder.Size() == 0 {
		return responseRecorder.WriteResponse()
	}

	decoded, err := m.getDecodedBufferFromResponse(&responseRecorder)
	if err != nil {
		m.logger.Error("error getting initial response", zap.Error(err))
		return responseRecorder.WriteResponse()
	}

	// Extract form request
	if err := r.ParseForm(); err != nil {
		return errors.Join(errors.New("failed to parse form"), err)
	}

	// Remove unsupported query parameters
	filterForm(&r.Form)

	// Return if no parameters remains
	if len(r.Form) == 0 {
		return responseRecorder.WriteResponse()
	}

	// Send to security middleware if defined
	if m.Security != nil {
		if err := m.Security.ProcessRequestForm(&r.Form); err != nil {
			if errors.Is(err, BypassRequestError) {
				return responseRecorder.WriteResponse()
			}

			var abortRequestError *AbortRequestError
			if errors.As(err, &abortRequestError) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}

			return err
		}

		// Return initial image if no parameters remains
		if len(r.Form) == 0 {
			return responseRecorder.WriteResponse()
		}
	}

	// Generate specific ETag ig necessary
	processedEtag := getProcessedImageEtag(responseRecorder.Header().Get("ETag"), &r.Form)
	if processedEtag != "" {
		responseRecorder.Header().Del("ETag") // Remove initial ETag
		w.Header().Set("ETag", processedEtag)

		// Check If-None-Match header to avoid reprocessing
		ifNoneMatchHeader := r.Header.Get("If-None-Match")
		if ifNoneMatchHeader != "" && ifNoneMatchHeader == processedEtag {
			w.WriteHeader(http.StatusNotModified)
			return nil
		}
	}

	// Parse options
	options, err := getOptions(&r.Form)
	if err != nil {
		m.logger.Error("error parsing options", zap.Error(err))
		return responseRecorder.WriteResponse()
	}

	newImage, err := bimg.NewImage(decoded).Process(options)
	if err != nil {
		m.logger.Error("error processing image", zap.Error(err))
		if m.OnFail == OnFailBypass {
			return responseRecorder.WriteResponse()
		}
		if m.OnFail == OnFailAbort {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		return err

	}

	// Remove proxied invalid header
	w.Header().Del("Content-Type")
	w.Header().Del("Content-Length")
	w.Header().Del("Content-Encoding")
	w.Header().Del("Vary")
	//w.Header().Del("ETag")

	// Set new headers
	w.Header().Set("Content-Length", strconv.Itoa(len(newImage)))
	w.Header().Set("Content-Type", "image/"+bimg.DetermineImageTypeName(newImage))

	if _, err = w.Write(newImage); err != nil {
		m.logger.Error("error writing processed image", zap.Error(err))
		if m.OnFail == OnFailBypass {
			return responseRecorder.WriteResponse()
		}
		if m.OnFail == OnFailAbort {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		return err
	}

	return nil
}

func (m *Middleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.NextBlock(0) {
			switch d.Val() {
			case "on_fail":
				// Check if argument provided
				if !d.NextArg() {
					return d.ArgErr()
				}
				m.OnFail = OnFail(d.Val())

				// Ensure there are no more arguments
				if d.NextArg() {
					return d.ArgErr() // More than one argument provided
				}

				break
			case "security":
				m.Security = &SecurityOptions{}
				if err := m.Security.UnmarshalCaddyfile(d); err != nil {
					return err
				}
				break

			default:
				return d.Errf("unexpected directive '%s' in image_processor block", d.Val())
			}
		}
	}
	return nil
}

func (m *Middleware) getDecodedBufferFromResponse(r *caddyhttp.ResponseRecorder) ([]byte, error) {

	encoding := (*r).Header().Get("Content-Encoding")
	if encoding == "" {
		return (*r).Buffer().Bytes(), nil
	}

	if encoding == "gzip" {
		decoder, err := gzip.NewReader((*r).Buffer())
		if err != nil {
			return nil, err
		}
		defer func(decoder *gzip.Reader) {
			err := decoder.Close()
			if err != nil {
				return
			}
		}(decoder)

		decodedOut := bytes.Buffer{}
		_, err = io.Copy(&decodedOut, decoder)
		if err != nil {
			return nil, err
		}
		return decodedOut.Bytes(), nil
	}

	if encoding == "zstd" {
		// Try decode zstd
		var decoder, err = zstd.NewReader((*r).Buffer(), zstd.WithDecoderConcurrency(0))
		if err != nil {
			return nil, err
		}
		defer decoder.Close()
		decodedOut := bytes.Buffer{}
		_, err = io.Copy(&decodedOut, decoder)
		if err != nil {
			return nil, err
		}
		return decodedOut.Bytes(), nil
	}

	return nil, fmt.Errorf("unsupported encoding: %s", encoding)

}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Middleware
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return &m, err
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddy.Validator             = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
	_ caddyfile.Unmarshaler       = (*Middleware)(nil)
)
