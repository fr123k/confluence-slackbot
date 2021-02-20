.PHONY: build

go-init:
	go mod init github.com/fr123k/confluence-slackbot
	go mod vendor

build:
	go build -o build/main cmd/main.go
	go test -v --cover ./...

run: build
	./build/main

clean:
	rm -rfv ./build
