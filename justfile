test:
    go test -v ./...

deps:
    docker compose up

term:
    go run term/main.go
