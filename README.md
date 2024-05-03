[![Build Static Releases](https://github.com/quix-labs/caddy-image-processor/actions/workflows/build-on-release.yml/badge.svg)](https://github.com/quix-labs/caddy-image-processor/actions/workflows/build-on-release.yml)

# Caddy Image Processor

This repository contains a CaddyServer module for processing images on the fly using libvips.

## Features

- Automatic image processing based on URL query parameters
- Supports resizing, rotating, cropping, quality adjustments, format conversion, and more
- Efficient processing using libvips

## Prerequisites

- [Caddy](https://caddyserver.com/) installed on your system
- [libvips](https://libvips.github.io/libvips/install.html) installed on your system
- [libvips-dev](https://libvips.github.io/libvips/install.html) installed on your system

## Installation and Configuration

### Using Docker

- Pull the Docker image from the GitHub Container Registry:
    ```bash
    docker pull ghcr.io/quix-labs/caddy-image-processor:latest
    ```

### Using xcaddy

- Before building the module, ensure you have `xcaddy` installed on your system. You can install it using the following
  command:

  ```bash
  go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest
  ```

- To build this module into Caddy, run the following command:

  ```bash
  CGO_ENABLED=1 xcaddy build --with github.com/quix-labs/caddy-image-processor
  ```

  This command compiles Caddy with the image processing module included.

### Using prebuilt assets

- You can also install the tool using release assets.

  Download the appropriate package from
  the [Releases page](https://github.com/quix-labs/caddy-image-processor/releases), and then follow the instructions
  provided for your specific platform.

## Usage

### Using Docker

```bash
docker run -p 80:80 -v $PWD/Caddyfile:/etc/caddy/Caddyfile -d ghcr.io/quix-labs/caddy-image-processor:latest
```

Your can see more information in the [official docker documentation for caddy](https://hub.docker.com/_/caddy)

### Using xcaddy build / prebuilt assets

```bash
/path/to/your/caddy run --config /etc/caddy/Caddyfile
```

Your can see more information in
the [official documentation for caddy](https://caddyserver.com/docs/build#package-support-files-for-custom-builds-for-debianubunturaspbian)

## Example Caddyfile

```plaintext
{
    order image_processor before respond
}

localhost {
    root * /your-images-directory
    file_server
    image_processor
}
```

In this example, all requests undergo processing by the image processor module before being served by the
caddy.

## Available Query Parameters

| Param | Name          | Description                                                                                             | Type                          |
|-------|---------------|---------------------------------------------------------------------------------------------------------|-------------------------------|
| h     | Height        | Image height                                                                                            | Integer                       |
| w     | Width         | Image width                                                                                             | Integer                       |
| ah    | AreaHeight    | Area height                                                                                             | Integer                       |
| aw    | AreaWidth     | Area width                                                                                              | Integer                       |
| t     | Top           | Y-coordinate of the top-left corner                                                                     | Integer                       |
| l     | Left          | X-coordinate of the top-left corner                                                                     | Integer                       |
| q     | Quality       | Image quality (JPEG compression)                                                                        | Integer (default 75)          |
| cp    | Compression   | Compression level (0-9, 0 = lossless)                                                                   | Integer                       |
| z     | Zoom          | Zoom level                                                                                              | Integer                       |
| crop  | Crop          | Whether cropping is enabled                                                                             | Boolean                       |
| en    | Enlarge       | Whether enlargement is enabled                                                                          | Boolean                       |
| em    | Embed         | Whether embedding is enabled                                                                            | Boolean                       |
| flip  | Flip          | Whether vertical flipping is enabled                                                                    | Boolean                       |
| flop  | Flop          | Whether horizontal flipping is enabled                                                                  | Boolean                       |
| force | Force         | Whether to force action                                                                                 | Boolean                       |
| nar   | NoAutoRotate  | Whether auto-rotation is disabled                                                                       | Boolean                       |
| np    | NoProfile     | Whether profile is disabled                                                                             | Boolean                       |
| itl   | Interlace     | Whether interlacing is enabled                                                                          | Boolean (default true)        |
| smd   | StripMetadata | Whether to strip metadata                                                                               | Boolean (default true)        |
| tr    | Trim          | Whether trimming is enabled                                                                             | Boolean                       |
| ll    | Lossless      | Whether compression is lossless                                                                         | Boolean                       |
| th    | Threshold     | Color threshold                                                                                         | Float                         |
| g     | Gamma         | Gamma correction                                                                                        | Float                         |
| br    | Brightness    | Brightness                                                                                              | Float                         |
| c     | Contrast      | Contrast                                                                                                | Float                         |
| r     | Rotate        | Rotation angle (45, 90, 135, 180, 235, 270, 315)                                                        | Integer                       |
| b     | GaussianBlur  | Gaussian blur level                                                                                     | Integer                       |
| bg    | Background    | Background color (white, black, red, magenta, blue, cyan, green, yellow, or hexadecimal format #RRGGBB) | Color                         |
| fm    | Type          | Image type (jpg, png, gif, webp, avif)                                                                  | Image Type (default original) |

## Planned Features

The following features are planned for future implementation:

- Sharp compliance: fit,...

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

## Credits

- [COLANT Alan](https://github.com/alancolant)
- [All Contributors](../../contributors)

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
