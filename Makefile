.PHONY: build run dev test vet tidy clean

build:
	cd backend && go build ./...

run:
	cd backend && go run cmd/main.go

dev:
	cd backend && air

test:
	cd backend && go test ./...

vet:
	cd backend && go vet ./...

tidy:
	cd backend && go mod tidy

clean:
	lsof -ti:8080 | xargs kill -9 2>/dev/null || true
