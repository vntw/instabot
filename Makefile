BINARY = instabot
GOARCH = amd64
MAIN = main/main.go
BUILD_DIR = build

.PHONY: all
all:
	$(MAKE) clean
	$(MAKE) linux
	$(MAKE) darwin

.PHONY: linux
linux:
	GOOS=linux GOARCH=${GOARCH} go build -o ${BUILD_DIR}/${BINARY}-linux-${GOARCH} ${MAIN_DIR}

.PHONY: darwin
darwin:
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BUILD_DIR}/${BINARY}-darwin-${GOARCH} ${MAIN_DIR}

.PHONY: clean
clean:
	-rm -r build/; \
	mkdir build/
