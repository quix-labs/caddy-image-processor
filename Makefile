# xcaddy parameters
XCADDY=xcaddy

# Main module path
MODULE_PATH=github.com/quix-labs/caddy-image-processor

# Output directory
OUT_DIR=out

# Binary names
BINARY_NAME=caddy

# Build targets
all: clean test build
build:
	$(XCADDY) build --output $(OUT_DIR)/$(BINARY_NAME) --with $(MODULE_PATH)=./
	chmod u+x $(OUT_DIR)/$(BINARY_NAME)
test:
	go test -v ./...
clean:
	rm -rf $(OUT_DIR)
run:
	$(XCADDY) run
	./$(OUT_DIR)/$(BINARY_NAME)
