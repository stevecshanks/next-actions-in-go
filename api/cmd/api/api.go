package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/stevecshanks/next-actions-in-go/api/internal/config"
	"github.com/stevecshanks/next-actions-in-go/api/internal/nextactions"
	"github.com/stevecshanks/next-actions-in-go/api/internal/trello"
)

func handleError(w http.ResponseWriter, err error) {
	fmt.Printf("Error: %s\n", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func actions(w http.ResponseWriter, req *http.Request) {
	cfg, err := config.FromEnvironment()
	if err != nil {
		handleError(w, err)
	}
	client := trello.Client{
		BaseURL: cfg.TrelloBaseURL,
		Key:     cfg.TrelloKey,
		Token:   cfg.TrelloToken,
	}

	fetcher := nextactions.Fetcher{Client: &client, Config: cfg}

	startTime := time.Now()

	actions, err := fetcher.Fetch()
	if err != nil {
		handleError(w, err)
	}

	fmt.Printf("Finished API requests, took %s\n", time.Since(startTime))

	if err := json.NewEncoder(w).Encode(map[string][]nextactions.Action{"data": actions}); err != nil {
		handleError(w, err)
	}
}

func main() {
	http.HandleFunc("/actions", actions)

	fmt.Printf("Listening on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
