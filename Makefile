
CMDBIN = gosn

run:
	go run cmd/main.go

build:
	go build -o bin/${CMDBIN} cmd/main.go
