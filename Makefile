GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
TARGET := obscure

all: $(TARGET)

clean: 
	rm bin/$(GOOS)/$(GOARCH)/$(TARGET)

deps:
	. gvp && gpm install

$(TARGET):
	. gvp && GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$(GOOS)/$(GOARCH)/$(TARGET) ./cmd/$(TARGET)

test: *.go **/*.go
	. gvp && go test ./...

install: $(TARGET) test
	cp bin/$(GOOS)/$(GOARCH)/$(TARGET) /usr/local/bin

.PHONY: clean deps $(TARGET)
