all: app

app: *.go go.mod go.sum
	go build -o app

.PHONY: lint
lint: fmt
	@go vet ./...
	@staticcheck ./...

.PHONY: fmt
fmt: install
	@goimports -l -w .

.PHONY: install
install:
	@go install golang.org/x/tools/cmd/goimports@v0.14.0
	@go install honnef.co/go/tools/cmd/staticcheck@v0.4.6
