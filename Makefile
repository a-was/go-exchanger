include .env
export $(shell sed 's/=.*//' .env)

run:
	CGO_ENABLED=0 go run .

test:
	go test ./...
