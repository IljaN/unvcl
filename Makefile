VERSION := $(shell git describe --tags --dirty --always)
LDFLAGS = -ldflags "-X main.Version=${VERSION}"

build:
	CGO_ENABLED=0 go build ${LDFLAGS} .

build_all: build/unvcl-amd64-linux build/unvcl-amd64-darwin build/unvcl-amd64-windows.exe build/unvcl-arm64-darwin

build/unvcl-amd64-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o build/unvcl-amd64-linux-${VERSION} main.go

build/unvcl-amd64-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o build/unvcl-amd64-darwin-${VERSION} main.go

build/unvcl-amd64-windows.exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o build/unvcl-amd64-windows-${VERSION}.exe main.go

build/unvcl-arm64-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o build/unvcl-arm64-darwin-${VERSION} main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: test-coverage
coverage:
	go test -coverprofile cover.out -v ./...

.PHONY: show-coverage
show-coverage: coverage
	go tool cover -html=cover.out

.PHONY: clean
clean:
	rm -rf ./build unvcl cover.out
