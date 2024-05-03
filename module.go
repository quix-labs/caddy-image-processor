package CADDY_FILE_SERVER

import (
	"bytes"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/h2non/bimg"
	"net/http"
	"strconv"
)

func init() {
	caddy.RegisterModule(Middleware{})
	httpcaddyfile.RegisterHandlerDirective("image_processor", parseCaddyfile)
	httpcaddyfile.RegisterDirectiveOrder("image_processor", "before", "respond")
}

type Middleware struct{}

func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.image_processor",
		New: func() caddy.Module { return new(Middleware) },
	}
}

func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

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

	options, err := getOptions(r)
	if err != nil {
		return err
	}

	newImage, err := bimg.NewImage(responseRecorder.Buffer().Bytes()).Process(options)
	if err != nil {
		return responseRecorder.WriteResponse()
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(newImage)))
	w.Header().Set("Content-Type", "image/"+bimg.DetermineImageTypeName(newImage))

	if _, err = w.Write(newImage); err != nil {
		return err
	}

	return nil
}

func (m *Middleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil

}

func getOptions(r *http.Request) (bimg.Options, error) {

	options := bimg.Options{
		Interlace:     true,
		StripMetadata: true,
	}

	type CustomProcessor struct {
		Func func(value string) error
	}
	parameters := map[string]interface{}{
		"h":     &options.Height,             // int
		"w":     &options.Width,              // int
		"ah":    &options.AreaHeight,         // int
		"aw":    &options.AreaWidth,          // int
		"t":     &options.Top,                // int
		"l":     &options.Left,               // int
		"q":     &options.Quality,            // int
		"cp":    &options.Compression,        // int
		"z":     &options.Zoom,               // int
		"crop":  &options.Crop,               // bool
		"en":    &options.Enlarge,            // bool
		"em":    &options.Embed,              // bool
		"flip":  &options.Flip,               // bool
		"flop":  &options.Flop,               // bool
		"force": &options.Force,              // bool
		"nar":   &options.NoAutoRotate,       // bool
		"np":    &options.NoProfile,          // bool
		"itl":   &options.Interlace,          // bool
		"smd":   &options.StripMetadata,      // bool
		"tr":    &options.Trim,               // bool
		"ll":    &options.Lossless,           // bool
		"th":    &options.Threshold,          // float64
		"g":     &options.Gamma,              // float64
		"br":    &options.Brightness,         // float64
		"c":     &options.Contrast,           // float64
		"r":     &options.Rotate,             // bimg.Angle
		"b":     &options.GaussianBlur.Sigma, // int
		"bg":    &options.Background,         // bimg.Color
		"fm":    &options.Type,               // bimg.Type
	}

	if err := r.ParseForm(); err != nil {
		return options, err
	}

	for param, _ := range r.Form {
		value := r.FormValue(param)
		if value == "" {
			continue
		}
		dest, exists := parameters[param]
		if !exists {
			continue
		}

		var err error
		switch dest.(type) {
		case *int:
			dest := dest.(*int)
			if *dest, err = strconv.Atoi(value); err != nil {
				return options, err
			}

		case *bool:
			dest := dest.(*bool)
			if *dest, err = strconv.ParseBool(value); err != nil {
				return options, err
			}

		case *float64:
			dest := dest.(*float64)
			if *dest, err = strconv.ParseFloat(value, 64); err != nil {
				return options, err
			}

		case *string:
			dest := dest.(*string)
			*dest = value

		case *bimg.Color:
			dest := dest.(*bimg.Color)

			if value == "white" {
				*dest = bimg.Color{255, 255, 255}
				break
			}
			if value == "black" {
				*dest = bimg.Color{0, 0, 0}
				break
			}
			if value == "red" {
				*dest = bimg.Color{255, 0, 0}
				break
			}
			if value == "magenta" {
				*dest = bimg.Color{255, 0, 255}
				break
			}
			if value == "blue" {
				*dest = bimg.Color{0, 0, 255}
				break
			}
			if value == "cyan" {
				*dest = bimg.Color{0, 255, 255}
				break
			}
			if value == "green" {
				*dest = bimg.Color{0, 255, 0}
				break
			}
			if value == "yellow" {
				*dest = bimg.Color{255, 255, 0}
				break
			}

			c := bimg.Color{}
			_, err := fmt.Sscanf(value, "#%02x%02x%02x", &c.R, &c.G, &c.B)
			if err != nil {
				return options, fmt.Errorf("possible values for '%s' are white,black,red,magenta,blue,cyan,green,yellow or #xxxxx hex string", param)
			}

			*dest = c

		case *bimg.Angle:
			dest := dest.(*bimg.Angle)
			angle, err := strconv.Atoi(value)
			if err != nil {
				return options, err
			}

			switch angle {
			case 45, 90, 135, 180, 235, 270, 315:
				*dest = bimg.Angle(angle)
			default:
				return options, fmt.Errorf("possible values for '%s' are 45, 90, 135, 180, 235, 270, 315", param)
			}

		case *bimg.ImageType:
			dest := dest.(*bimg.ImageType)
			switch value {
			case "jpg", "jpeg":
				*dest = bimg.JPEG
			case "png":
				*dest = bimg.PNG
			case "gif":
				*dest = bimg.GIF
			case "webp":
				*dest = bimg.WEBP
			case "avif":
				*dest = bimg.AVIF
			default:
				return options, fmt.Errorf("possible values for '%s' are jpg, jpeg, png, gif, webp, avif", param)
			}
		}
	}

	return options, nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Middleware
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

// Interface guards
var (
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
	_ caddyfile.Unmarshaler       = (*Middleware)(nil)
)
