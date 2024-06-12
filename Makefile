gen-doc:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --pd -g pkg/api/http/rate-server.go

	swag fmt
