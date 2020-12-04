package main

import (
	"github.com/pirsch-analytics/go-sdk"
	"log"
	"net/http"
)

const (
	// Client ID, secret, and hostname for testing.
	// Replace them with your own.
	clientID     = "i9OulOrSI0b5EBNQbb8vg5zfEg2zhQ5q"
	clientSecret = "FgwcclcvpTK75oqmbc6UxmJTXX3iCK8JcqPvz7ozzJzouc39KLEeriXQy45Myq92"
	hostname     = "first.page"
)

func main() {
	log.Println("Visit http://localhost:1414")

	// Create a client for Pirsch.
	client := pirsch.NewClient(clientID, clientSecret, hostname, &pirsch.ClientConfig{
		BaseURL: "http://localhost.com:9999",
	})

	// Add a handler to serve a page.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// To analyze traffic, send page hits to Pirsch.
		// You can control what gets send to Pirsch, in case the page was not found for example.
		if r.URL.Path == "/" {
			if err := client.Hit(r); err != nil {
				log.Println(err)
			}
		}

		w.Write([]byte("<h1>Hello from Pirsch!</h1>"))
	})
	http.ListenAndServe(":1414", nil)
}
