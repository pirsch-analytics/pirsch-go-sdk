.PHONY: deps test

deps:
	go get -u -t ./...
	go mod tidy
	go mod vendor

fix:
	go fix ./...

test:
	go test -cover -race github.com/pirsch-analytics/pirsch-go-sdk/v2/pkg
