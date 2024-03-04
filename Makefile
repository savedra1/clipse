BINARY_NAME=clipse
INSTALL_DIR="/usr/bin/"
 
all: build run
 
install:
	go build -o "${INSTALL_DIR}${BINARY_NAME}"

build:
	go build -o "${BINARY_NAME}"

run:
	go build -o ${BINARY_NAME}
	./${BINARY_NAME}
 
clean:
	go clean
	rm ${BINARY_NAME}