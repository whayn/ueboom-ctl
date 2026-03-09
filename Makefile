BINARY_NAME=ueboom-ctl
MAIN_PATH=./cmd/ueboom-ctl

.PHONY: build install clean help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build    Build the binary"
	@echo "  install  Install the binary to GOBIN"
	@echo "  clean    Remove the binary"

build:
	go build -v -o $(BINARY_NAME) $(MAIN_PATH)

install:
	go install -v $(MAIN_PATH)

clean:
	rm -f $(BINARY_NAME)
