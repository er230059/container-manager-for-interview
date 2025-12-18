gen-doc:
	 swag init --parseDependency -g cmd/main.go -o docs

build:
	go mod tidy
	go build -o build/main cmd/main.go
