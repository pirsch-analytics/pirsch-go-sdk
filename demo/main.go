package main

import (
	"github.com/pirsch-analytics/pirsch-go-sdk"
	"log"
	"net/http"
	"os"
)

const (
	// Client ID, secret, and hostname for testing.
	// Replace them with your own.
	clientID     = ""
	clientSecret = ""
	hostname     = "pirsch.io"
)

func main() {
	log.Println("Visit http://localhost:1414")

	// Create a client for Pirsch.
	client := pirsch.NewClient(clientID, clientSecret, hostname, &pirsch.ClientConfig{
		Logger: log.New(os.Stdout, "", 0),
	})

	// Add a handler to serve a page.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// To analyze traffic, send page hits to Pirsch.
		// You can control what gets send to Pirsch, in case the page was not found for example.
		if r.URL.Path == "/" {
			if err := client.Hit(r); err != nil {
				log.Println(err)
			}

			log.Println("Hit!")
		}

		w.Write([]byte("<h1>Hello from Pirsch!</h1>"))
	})
	http.ListenAndServe(":1414", nil)
}
