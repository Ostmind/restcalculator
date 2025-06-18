IMAGE_NAME=api-calculator

apidoc:
	swag init -g cmd/restcalculator/main.go
fmt:
	go fmt ./...
lint:
	golangci-lint run
test:
	go test ./...
coverage:
	go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
build:
	docker build -t $(IMAGE_NAME) .