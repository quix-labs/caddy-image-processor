# Caddy Image Processor

This repository contains a CaddyServer module for processing images on the fly using libvips.

## Features
- Automatic image processing based on URL query parameters
- Supports resizing, rotating, cropping, quality adjustments, format conversion, and more
- Efficient processing using libvips

## Prerequisites
- [Caddy](https://caddyserver.com/) installed on your system
- [libvips](https://libvips.github.io/libvips/install.html) installed on your system

## Building with xcaddy

Before building the module, ensure you have `xcaddy` installed on your system. You can install it using the following command:

```bash
go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest
```

To build this module into Caddy, run the following command:

```bash
xcaddy build --with github.com/quix-labs/caddy-image-processor
```

This command compiles Caddy with the image processing module included.

## Usage

Follow these steps to utilize the image processing capabilities:

1. Install Caddy and libvips on your system.
2. Build Caddy with the image processing module using xcaddy.
3. Configure your Caddyfile to include the image processing module for specific routes or sites.
4. Start Caddy, and access your images with processing options via URL query parameters.

## Available Query Parameters

- or: Orientation (e.g., 90, 180, 270)
- crop: Crop (1 for true, 0 for false)
- w: Width
- h: Height
- blur: Blur amount
- q: Quality
- fm: Format (e.g., jpg, png, gif, webp, avif)

## Example Caddyfile
```plaintext
example.com {
  route /images* {
    reverse_proxy localhost:8080
    image_processor
  }
}
```

In this example, requests to `/images*` undergo processing by the image processor module before being served by the reverse proxy.

## Planned Features

The following features are planned for future implementation:

- FLIP parameter
- CROP NOT GLIDE COMPLIANT parameter adjustments
- Additional parameters: fit, dpr, bri, con, gam, sharp
- Parameters for adding watermark: pixel, filt, mark, markw, markh, markx, marky, markpad, markpos, markalpha, bg, border