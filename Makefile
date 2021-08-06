tidy:
	go mod tidy
run: tidy
	go run main.go
test:
	go test ./... -v -cover