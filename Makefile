BINARY_NAME=clipse
INSTALL_DIR?=$(HOME)/.local/bin

wayland:
	CGO_ENABLED=0 go build -tags wayland -o $(BINARY_NAME)

x11:
	go build -tags linux -o $(BINARY_NAME)

darwin:
	go build -tags darwin -o $(BINARY_NAME)

run: wayland
	./$(BINARY_NAME)

install: wayland
	install -Dm 755 $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

clean:
	go clean
	rm -f $(BINARY_NAME)

test:
	go test -v ./...
	