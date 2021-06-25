package main

import (
	"github.com/pirsch-analytics/pirsch-go-sdk"
	"log"
	"net/http"
	"os"
)

const (
	// Client ID, secret, and hostname.
	// Replace them with your own.
	clientID     = "VPvOChTcKhn8gz0Xni0TaKY4C0PkuyKP"
	clientSecret = "MyYPXtEHKoGVZHNFuVDyBEMYsKZWvKrsPsSyHYqC4oRS2gyv62a0WRiEei4AryAE"
	hostname     = "pirsch.io"
)

func main() {
	log.Println("Visit http://localhost:1414")

	// Create a client for Pirsch.
	client := pirsch.NewClient(clientID, clientSecret, hostname, &pirsch.ClientConfig{
		Logger:  log.New(os.Stdout, "", 0),
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

			log.Println("Hit!")
		}

		w.Write([]byte("<h1>Hello from Pirsch!</h1>"))
	})

	// Add a handler to read statistics.
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		// Read the domain for this client.
		domain, err := client.Domain()

		if err == nil {
			w.Write([]byte("<p>Hostname: " + domain.Hostname + "</p>"))
		} else {
			w.Write([]byte(err.Error()))
			return
		}
	})

	http.ListenAndServe(":1414", nil)
}
