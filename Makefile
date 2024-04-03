.PHONY: build test clean

APP_NAME := pca 
VERSION := 0.0.3

build:
	cd cmd/$(APP_NAME) && go build -o ../../bin/$(APP_NAME) -ldflags "-X main.AppVersion=$(VERSION)" \
	&& cd .. && \
	if ! git rev-parse -q --verify v$(VERSION); then \
		git tag v$(VERSION) -m release; \
	else \
		echo "Tag v$(VERSION) already exists."; \
	fi

test:
	go test ./pkg/...

clean:
	rm -f bin/$(APP_NAME)
