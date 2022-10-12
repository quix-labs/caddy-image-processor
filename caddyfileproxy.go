package CADDY_FILE_SERVER

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/h2non/bimg"
	//_ "github.com/h2non/bimg"
	"net/http"
	"strconv"
)

func init() {
	caddy.RegisterModule(ProxyMiddleware{})
}

type ProxyMiddleware struct {
}

func (ProxyMiddleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.image_proxy",
		New: func() caddy.Module { return new(ProxyMiddleware) },
	}
}

func (m ProxyMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	//Automatic return if not options set
	if r.URL.RawQuery == "" {
		return next.ServeHTTP(w, r)
	}

	rw := &ResponseWrapper{w: w}
	err := next.ServeHTTP(rw, r)
	if err != nil {
		return err
	}

	options, err := getOptions(r)
	if err != nil {
		return err
	}

	fmt.Print(options)

	_, err = rw.Done()
	if err != nil {
		return err
	}
	return nil
}

func getOptions(r *http.Request) (bimg.Options, error) {
	options := bimg.Options{
		Enlarge:   true,
		Interlace: true,
	}

	if len(r.FormValue("or")) > 0 {
		orientation, err := strconv.Atoi(r.FormValue("or"))
		if err != nil {
			return options, err
		}
		options.Rotate = bimg.Angle(orientation)
	}
	//@TODO FLIP

	//@TODO CROP NOT GLIDE COMPLIANT
	if len(r.FormValue("crop")) > 0 {
		crop, err := strconv.Atoi(r.FormValue("crop"))
		if err != nil {
			return options, err
		}
		options.Crop = crop == 1
	}
	if len(r.FormValue("w")) > 0 {
		width, err := strconv.Atoi(r.FormValue("q"))
		if err != nil {
			return options, err
		}
		options.Width = width
	}
	if len(r.FormValue("h")) > 0 {
		height, err := strconv.Atoi(r.FormValue("h"))
		if err != nil {
			return options, err
		}
		options.Height = height
	}
	//@TODO fit,dpr,bri,con,gam,sharp,
	if len(r.FormValue("blur")) > 0 {
		blur, err := strconv.ParseFloat(r.FormValue("blur"), 10)
		if err != nil {
			return options, err
		}
		options.GaussianBlur.Sigma = blur
	}
	//@TODO pixel,filt,mark,markw,markh,markx,marky,markpad,markpos,markalpha,bg,border

	if len(r.FormValue("q")) > 0 {
		quality, err := strconv.Atoi(r.FormValue("q"))
		if err != nil {
			return options, err
		}
		options.Quality = quality
	}

	if len(r.FormValue("fm")) > 0 {
		format := r.FormValue("fm")
		switch format {
		case "jpg", "jpeg":
			options.Type = bimg.JPEG
		case "png":
			options.Type = bimg.PNG
		case "gif":
			options.Type = bimg.GIF
		case "webp":
			options.Type = bimg.WEBP
		case "avif":
			options.Type = bimg.AVIF
		}
	}

	return options, nil
}

// Interface guards
var (
	_ caddyhttp.MiddlewareHandler = (*ProxyMiddleware)(nil)
)
