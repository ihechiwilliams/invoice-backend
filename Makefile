env:
	cp -v .env.example .env

local:
	go run cmd/server/*.go

# Run go generate locally without docker container
generate:
	go run github.com/vektra/mockery/v2@v2.43.0
	go generate ./...