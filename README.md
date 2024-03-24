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
CGO_ENABLED=1 xcaddy build --with github.com/quix-labs/caddy-image-processor
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

## Development

To contribute to the development of Caddy Image Processor, follow these steps:

1. Make sure you have Go installed on your system.
2. Clone this repository to your local machine:
   ```bash
   git clone https://github.com/quix-labs/caddy-image-processor.git
   ```
   
3. Navigate to the project directory:
4. Install `xcaddy` if you haven't already:
    ```bash
    go get -u github.com/caddyserver/xcaddy/cmd/xcaddy
    ```
5. Make your changes in the source code.
6. Run tests to ensure your changes haven't introduced any issues:
    ```bash
   make test
    ```
7. If tests pass, you can build the project:
    ```bash
   make build
    ```
8. To run the project in development mode, use the following command:
    ```bash
   make run
    ```
9. Once you're satisfied with your changes, create a pull request to the main branch of the repository for review.