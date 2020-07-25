.PHONY: watch-test watch-run fetch-test-data
watch-test:
	fswatch -o . | xargs -n1 -I{} go test -v ./...

watch-run:
	fswatch -o . | xargs -n1 -I{} go run main.go

fetch-test-data:
	go run main.go -command="fetch-repo-metadata-list-for-user" -user="pfeilbr" -output="testdata/repos/pfeilbr/repo-list.json"
