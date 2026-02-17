# Pirsch Golang SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/pirsch-analytics/pirsch-go-sdk/v2?status.svg)](https://pkg.go.dev/github.com/pirsch-analytics/pirsch-go-sdk/v2?status)

This is the official Golang client SDK for Pirsch. For details, please check out our [documentation](https://docs.pirsch.io/).

## Install

```
go get github.com/pirsch-analytics/pirsch-go-sdk/v2
```

## Usage

```go
package main

import (
	"log"

	pirsch "github.com/pirsch-analytics/pirsch-go-sdk/v2/pkg"
)

func main() {
	// Create a new client using the client ID and secret you've created on the dashboard.
	client := pirsch.NewClient("client_id", "client_secret", nil)

	// Get the dashboard domain for this client (you should handle the error).
	domain, _ := client.Domain()
	
	// Print the hostname for the dashboard.
	log.Println(domain.Hostname)
}
```

## Changelog

See [CHANGELOG.md](CHANGELOG.md).

## License

MIT
