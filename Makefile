.PHONY: watch-test watch-run
watch-test:
	fswatch -o . | xargs -n1 -I{} go test -v ./...

watch-run:
	fswatch -o . | xargs -n1 -I{} go run main.go
