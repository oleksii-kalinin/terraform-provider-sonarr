package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/oleksii-kalinin/terraform-provider-sonarr/pkg/sonarr"
)

// main initializes environment and Sentry, validates required Sonarr environment variables,
// and exercises the Sonarr client by fetching system status, retrieving a series (ID 78),
// creating a series, waiting 10 seconds, and then deleting the created series.
// Errors from Sonarr operations are reported to Sentry and logged; missing required
// environment variables or a failed Sentry initialization cause the program to exit.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error parsing .env: %s", err)
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	apiKey := os.Getenv("SONARR_API_KEY")
	if apiKey == "" {
		message := "SONARR_API_KEY environment variable should be set"
		sentry.CaptureException(errors.New(message))
		log.Fatal(message)
	}

	url := os.Getenv("SONARR_URL")
	if url == "" {
		message := "SONARR_URL environment variable should be set"
		sentry.CaptureException(errors.New(message))
		log.Fatal(message)
	}

	client := sonarr.NewClient(url, apiKey)

	status, err := client.GetSystemStatus()
	if err != nil {
		sentry.CaptureException(err)
		log.Println(err)
		return
	}
	fmt.Println(status)

	series, err := client.GetSeries(78)
	if err != nil {
		sentry.CaptureException(err)
		log.Println(err)
		return
	}
	fmt.Println(series)

	wantedShow := sonarr.Series{
		TvdbID:           289590,
		Title:            "Mr. Robot",
		Monitored:        false,
		RootFolderPath:   "/media/series",
		QualityProfileId: 1,
	}

	wanted, err := client.CreateSeries(&wantedShow)
	if err != nil {
		sentry.CaptureException(err)
		log.Println(err)
		return
	}
	fmt.Println(wanted)
	time.Sleep(10 * time.Second)

	err = client.DeleteSeries(wanted.Id, true)
	if err != nil {
		sentry.CaptureException(err)
		log.Println(err)
		return
	}
}