.PHONY: build test clean

APP_NAME=golv
VERSION=0.0.1

build:
	cd cmd/goLV && go build -o ../../bin/$(APP_NAME) -ldflags "-X main.AppVersion=$(VERSION)"

test:
	go test ./pkg/...

clean:
	rm -f bin/$(APP_NAME)
