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

# reads the tags markdown table tags.md and print out the description and then the tag
# this printpout was used to extract the constants in the pkg/nlp/tag.go file.
extract-tags:
	tail -n +3 tags.md | sed 's/`/ /g' | awk '{split($$0,a,"|"); print a[3],a[2]}'
