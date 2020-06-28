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

type APIError struct {
	Detail string `json:"detail"`
}

func handleError(w http.ResponseWriter, err error) {
	fmt.Printf("Error: %s\n", err)

	apiErrors := []APIError{{err.Error()}}

	body, err := json.Marshal(map[string][]APIError{"errors": apiErrors})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(body)
	if err != nil {
		panic(err)
	}
}

func actions(w http.ResponseWriter, req *http.Request) {
	cfg, err := config.FromEnvironment()
	if err != nil {
		handleError(w, err)
	}
	client := trello.Client{
		Key:   cfg.TrelloKey,
		Token: cfg.TrelloToken,
	}

	fetcher := nextactions.Fetcher{Client: &client, Config: cfg}

	startTime := time.Now()

	actions, err := fetcher.Fetch()
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Printf("Finished API requests, took %s\n", time.Since(startTime))

	if err := json.NewEncoder(w).Encode(map[string][]nextactions.Action{"data": actions}); err != nil {
		handleError(w, err)
	}
}

func main() {
	http.HandleFunc("/actions", actions)

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
