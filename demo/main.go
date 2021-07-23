package main

import (
	"encoding/json"
	"github.com/pirsch-analytics/pirsch-go-sdk"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// Client ID, secret, and hostname.
	// Replace them with your own.
	clientID     = "2zH9LVKwEHv9nCK8Nr81HHeLF9Olz2ip"
	clientSecret = "8BhHk3GBvUdgWhdQygqzxRDZKTYadh8URARmSFxqlhDiWPPRT0ycWOi7kdejpZHY"
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

	// Add a handler to send an event.
	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		if err := client.Event("My First Event", 42, map[string]string{"hello": "world"}, r); err != nil {
			log.Println(err)
		}

		w.Write([]byte("<h1>Event sent!</h1>"))
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

		// Read statistics using the domain ID from above.
		filter := &pirsch.Filter{
			DomainID: domain.ID,
			From:     time.Now().Add(-time.Hour * 24 * 7), // one week
			To:       time.Now(),
		}
		visitors, _ := client.Visitors(filter)
		w.Write([]byte("<h2>Visitors</h2><pre>"))
		visitorsJson, _ := json.Marshal(visitors)
		w.Write(visitorsJson)
		w.Write([]byte("</pre>"))

		pages, _ := client.Pages(filter)
		w.Write([]byte("<h2>Pages</h2><pre>"))
		pagesJson, _ := json.Marshal(pages)
		w.Write(pagesJson)
		w.Write([]byte("</pre>"))

		sessionDuration, _ := client.SessionDuration(filter)
		w.Write([]byte("<h2>Session Duration</h2><pre>"))
		sessionDurationJson, _ := json.Marshal(sessionDuration)
		w.Write(sessionDurationJson)
		w.Write([]byte("</pre>"))

		timeOnPage, _ := client.TimeOnPage(filter)
		w.Write([]byte("<h2>Time on Page</h2><pre>"))
		timeOnPageJson, _ := json.Marshal(timeOnPage)
		w.Write(timeOnPageJson)
		w.Write([]byte("</pre>"))

		utmSource, _ := client.UTMSource(filter)
		w.Write([]byte("<h2>UTM Source</h2><pre>"))
		utmSourceJson, _ := json.Marshal(utmSource)
		w.Write(utmSourceJson)
		w.Write([]byte("</pre>"))

		utmMedium, _ := client.UTMMedium(filter)
		w.Write([]byte("<h2>UTM Medium</h2><pre>"))
		utmMediumJson, _ := json.Marshal(utmMedium)
		w.Write(utmMediumJson)
		w.Write([]byte("</pre>"))

		utmCampaign, _ := client.UTMCampaign(filter)
		w.Write([]byte("<h2>UTM Campaign</h2><pre>"))
		utmCampaignJson, _ := json.Marshal(utmCampaign)
		w.Write(utmCampaignJson)
		w.Write([]byte("</pre>"))

		utmContent, _ := client.UTMContent(filter)
		w.Write([]byte("<h2>UTM Content</h2><pre>"))
		utmContentJson, _ := json.Marshal(utmContent)
		w.Write(utmContentJson)
		w.Write([]byte("</pre>"))

		utmTerm, _ := client.UTMTerm(filter)
		w.Write([]byte("<h2>UTM Term</h2><pre>"))
		utmTermJson, _ := json.Marshal(utmTerm)
		w.Write(utmTermJson)
		w.Write([]byte("</pre>"))

		conversionGoals, _ := client.ConversionGoals(filter)
		w.Write([]byte("<h2>Conversion Goals</h2><pre>"))
		conversionGoalsJson, _ := json.Marshal(conversionGoals)
		w.Write(conversionGoalsJson)
		w.Write([]byte("</pre>"))

		events, _ := client.Events(filter)
		w.Write([]byte("<h2>Events</h2><pre>"))
		eventsJson, _ := json.Marshal(events)
		w.Write(eventsJson)
		w.Write([]byte("</pre>"))

		filter.Event = "My First Event"
		filter.EventMetaKey = "hello"
		eventMetadata, err := client.EventMetadata(filter)
		w.Write([]byte("<h2>Event Metadata</h2><pre>"))
		eventMetadataJson, err := json.Marshal(eventMetadata)
		w.Write(eventMetadataJson)
		w.Write([]byte("</pre>"))
		filter.Event = ""
		filter.EventMetaKey = ""

		growth, _ := client.Growth(filter)
		w.Write([]byte("<h2>Growth</h2><pre>"))
		growthJson, _ := json.Marshal(growth)
		w.Write(growthJson)
		w.Write([]byte("</pre>"))

		activeVisitors, _ := client.ActiveVisitors(filter)
		w.Write([]byte("<h2>Active Visitors</h2><pre>"))
		activeVisitorsJson, _ := json.Marshal(activeVisitors)
		w.Write(activeVisitorsJson)
		w.Write([]byte("</pre>"))

		timeOfDay, _ := client.TimeOfDay(filter)
		w.Write([]byte("<h2>Time of Day</h2><pre>"))
		timeOfDayJson, _ := json.Marshal(timeOfDay)
		w.Write(timeOfDayJson)
		w.Write([]byte("</pre>"))

		languages, _ := client.Languages(filter)
		w.Write([]byte("<h2>Languages</h2><pre>"))
		languagesJson, _ := json.Marshal(languages)
		w.Write(languagesJson)
		w.Write([]byte("</pre>"))

		referrer, _ := client.Referrer(filter)
		w.Write([]byte("<h2>Referrer</h2><pre>"))
		referrerJson, _ := json.Marshal(referrer)
		w.Write(referrerJson)
		w.Write([]byte("</pre>"))

		os, _ := client.OS(filter)
		w.Write([]byte("<h2>OS</h2><pre>"))
		osJson, _ := json.Marshal(os)
		w.Write(osJson)
		w.Write([]byte("</pre>"))

		browser, _ := client.Browser(filter)
		w.Write([]byte("<h2>Browser</h2><pre>"))
		browserJson, _ := json.Marshal(browser)
		w.Write(browserJson)
		w.Write([]byte("</pre>"))

		country, _ := client.Country(filter)
		w.Write([]byte("<h2>Country</h2><pre>"))
		countryJson, _ := json.Marshal(country)
		w.Write(countryJson)
		w.Write([]byte("</pre>"))

		platform, _ := client.Platform(filter)
		w.Write([]byte("<h2>Platform</h2><pre>"))
		platformJson, _ := json.Marshal(platform)
		w.Write(platformJson)
		w.Write([]byte("</pre>"))

		screen, _ := client.Screen(filter)
		w.Write([]byte("<h2>Screen</h2><pre>"))
		screenJson, _ := json.Marshal(screen)
		w.Write(screenJson)
		w.Write([]byte("</pre>"))

		keywords, _ := client.Keywords(filter)
		w.Write([]byte("<h2>keywords</h2><pre>"))
		keywordsJson, _ := json.Marshal(keywords)
		w.Write(keywordsJson)
		w.Write([]byte("</pre>"))
	})

	http.ListenAndServe(":1414", nil)
}
