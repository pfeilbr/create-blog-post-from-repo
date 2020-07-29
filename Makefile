.PHONY: watch-test watch-run fetch-test-data build
watch-test:
	fswatch -o *.go templates/*.* .env | xargs -n1 -I{} go test -v ./...

watch-run:
	fswatch -o . | xargs -n1 -I{} go run main.go

fetch-test-data:
	go run main.go -command="fetch-and-save-repos-for-user" -user="pfeilbr" -output="testdata/repos/pfeilbr/repo-list.json"

build:
	go build

install:
	cp create-blog-post-from-repo ~/bin
	chmod a+x ~/bin/create-blog-post-from-repo
