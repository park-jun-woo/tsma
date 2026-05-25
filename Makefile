VERSION := v0.2.2

.PHONY: install
install:
	go install -ldflags "-X main.Version=$(VERSION)" ./cmd/tsma

.PHONY: build
build:
	go build -ldflags "-X main.Version=$(VERSION)" -o tsma ./cmd/tsma

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -f tsma
