[![Build Static Releases](https://github.com/quix-labs/caddy-image-processor/actions/workflows/build-on-release.yml/badge.svg)](https://github.com/quix-labs/caddy-image-processor/actions/workflows/build-on-release.yml)
[![Build Docker](https://github.com/quix-labs/caddy-image-processor/actions/workflows/docker-on-release.yml/badge.svg)](https://github.com/quix-labs/caddy-image-processor/actions/workflows/docker-on-release.yml)

# Caddy Image Processor

This repository contains a CaddyServer module for processing images on the fly using libvips.

## Features

- Automatic image processing based on URL query parameters
- Supports resizing, rotating, cropping, quality adjustments, format conversion, and more
- Efficient processing using libvips

## Prerequisites

- [libvips](https://libvips.github.io/libvips/install.html) installed on your system
- [libvips-dev](https://libvips.github.io/libvips/install.html) installed on your system (only for xcaddy build)

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

### Using file_server

```plaintext
localhost {
    root /your-images-directory
    file_server
    image_processor
}
```

### Using reverse_proxy

```plaintext
localhost {
    reverse_proxy your-domain.com
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

## Examples

* Resize an image to a width of 400 pixels and a height of 300 pixels:
    * http://example.com/image.jpg?w=400&h=300
* Crop an image to a width of 200 pixels and a height of 150 pixels starting from the top-left corner (x=50, y=50):
    * http://example.com/image.jpg?w=200&h=150&t=50&l=50&crop=true
* Adjust the quality of the image to 80:
    * http://example.com/image.jpg?q=80
* Convert an image to PNG format and apply a Gaussian blur of 5:
    * http://example.com/image.jpg?fm=png&b=5
* Rotate an image by 180 degrees and flip it horizontally:
    * http://example.com/image.jpg?r=180&flop=true
* Apply a color threshold of 0.5 and adjust the brightness to -10:
    * http://example.com/image.jpg?th=0.5&br=-10
* Convert an image to AVIF format with lossless compression:
    * http://example.com/image.jpg?fm=avif&ll=true

## Advanced Configuration

This configuration allows you to control error handling with `on_fail` and `on_security_fail`.

You can also manage query parameter processing using `allowed_params` and `disallowed_params`.

This gives you fine-grained control over image processing in your Caddy server.


### Example with `on_fail` and Security Configuration
```plaintext
localhost:80 {
    root test-dataset
    file_server
	
    image_processor {
	    
        # Serve original image if image in unprocessable
        on_fail bypass	    
        
        # Return 500 Internal Server Error if processing fails
        # on_fail abort	    
        
        security {

            # Use ignore to remove param from processing, all valid param are processed
            on_security_fail ignore

            # Use abort to return 400 Bad Request when fails
            # on_security_fail abort

            # Use bypass to serve original image without processing
            # on_security_fail bypass

            # Explicitely disable rotate capabilities
            disallowed_params r
            
            # As an alternative use this to only accept width and height processing 
            # allowed_params w h 
            
            constraints {
                h range 60 480

                w {
                    values 60 130 240 480 637

                    # Shortcut range 60 637
                    range {
                        from 60
                        to 637
                    }
                }
            }
        }
    }
}
```

### Explanation:

* `on_fail`:
    * `bypass` (default value): If any error occurs, the original, unprocessed image will be returned.
    * `abort`: If an error occurs, a 500 Internal Server Error response will be returned.


* `on_security_fail`:
    * `ignore` (default value): If any security checks fail, they are ignored, and the image processing continues.
    * `bypass`: If any security checks fail, the original, unprocessed image will be returned.
    * `abort`: If any security checks fail, a 400 Bad Request response will be returned.


* **Security Configuration** (`disallowed_params` vs `allowed_params`):
  * `disallowed_params`: Specifies which query parameters are not allowed.
    
    For example, parameters like w (width) and r (rotation) can be restricted.

  * `allowed_params`: Specify which query parameters are allowed. As an alternative to `disallowed_params`.

  *  **Important**: You cannot use both allowed_params and disallowed_params in the same configuration.
  *  `constraints`: You san specify constraints for each parameter (see example)


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
