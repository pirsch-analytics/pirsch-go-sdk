.PHONY: deps test

deps:
	go get -u -t ./...

test:
	go test -cover -race github.com/pirsch-analytics/pirsch-go-sdk/v2/pkg
