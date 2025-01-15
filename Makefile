build-dummy:
	go build -o build/dummy cmd/dummy/main.go

run-proxy:
	go run cmd/proxy/main.go

run-client:
	go run cmd/httpClient/main.go