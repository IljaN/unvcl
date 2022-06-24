VERSION := $(shell git describe --tags --dirty --always)
LDFLAGS = -ldflags "-X main.Version=${VERSION}"

build:
	CGO_ENABLED=0 go build ${LDFLAGS} .

build_all: bin/unvcl-amd64-linux bin/unvcl-amd64-darwin bin/unvcl-amd64-windows.exe

bin/unvcl-amd64-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/unvcl-amd64-linux main.go

bin/unvcl-amd64-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/unvcl-amd64-darwin main.go

bin/unvcl-amd64-windows.exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/unvcl-amd64-windows.exe main.go

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
	rm -rf ./bin unvcl cover.out
