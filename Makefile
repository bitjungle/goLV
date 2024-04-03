.PHONY: build test clean

APP_NAME=pca
VERSION=0.0.2

build:
	cd cmd/$(APP_NAME) && go build -o ../../bin/$(APP_NAME) -ldflags "-X main.AppVersion=$(VERSION)"

test:
	go test ./pkg/...

clean:
	rm -f bin/$(APP_NAME)
