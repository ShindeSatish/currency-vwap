update:
	@echo "Updating Packages..."
	@go get -u
	@go mod tidy
run:
	@go run main.go

build:
	@go build -o currency-vwap main.go

test:
	@go test ./...
