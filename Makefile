.PHONY: swagger build run

swagger:
	swag init -g ./cmd/ibge-api/main.go

build:
	go build -o ./bin/ibge-service ./cmd/ibge-api

format:
# Verificar formatação
	gofmt -s -l .
# Corrigir formatação
	go fmt ./...

health:
	curl -f http://localhost:9080/health || exit 1
	
run:
	go run ./cmd/ibge-api/main.go

test:
	go test -v ./...

test-coverage:
	go test --coverprofile=coverage.out ./...

verify:
	go mod verify
	go mod tidy

docker-build:
	docker build -t ibge:latest .

docker-run:
	docker run --name ibge -p 9080:9080 -d ibge:latest

clean:
	rm -rf ./bin
	rm -rf ./coverage.out
	docker rm -f ibge || true
	