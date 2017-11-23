test:
	go test -v ./...
run:
	go run main.go        
build:
	docker run --rm -v $(HOME)/src/go:/go -w /go/src/github.com/maddevsio/nambataxi-telegram-stats-bot --name go golang:1.8 env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o statbot
