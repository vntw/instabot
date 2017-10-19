BINARY = instabot
GOARCH = amd64
MAIN_FILE = main/main.go
BUILD_DIR = build

.PHONY: all
all:
	$(MAKE) clean
	$(MAKE) linux
	$(MAKE) darwin

.PHONY: linux
linux:
	GOOS=linux GOARCH=${GOARCH} go build -o ${BUILD_DIR}/${BINARY}-linux-${GOARCH} ${MAIN_FILE}

.PHONY: darwin
darwin:
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BUILD_DIR}/${BINARY}-darwin-${GOARCH} ${MAIN_FILE}

.PHONY: clean
clean:
	-rm -r build/; \
	mkdir build/
