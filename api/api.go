package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"next-actions/api/trello"
)

// Action represents a "next action" in GTD
type Action struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// MarshalJSON adds fields required by JSON-API to an Action
func (a *Action) MarshalJSON() ([]byte, error) {
	type AliasedAction Action
	return json.Marshal(&struct {
		Type string `json:"type"`
		*AliasedAction
	}{
		Type:          "actions",
		AliasedAction: (*AliasedAction)(a),
	})
}

func handleError(w http.ResponseWriter, err error) {
	fmt.Printf("Error: %s\n", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func actions(w http.ResponseWriter, req *http.Request) {
	baseURL, err := url.Parse(os.Getenv("TRELLO_BASE_URL"))
	if err != nil {
		handleError(w, err)
	}
	client := trello.Client{BaseURL: baseURL, Key: os.Getenv("TRELLO_KEY"), Token: os.Getenv("TRELLO_TOKEN")}

	cards, err := client.ListOwnedCards()
	if err != nil {
		handleError(w, err)
	}

	actions := make([]Action, 0)
	for _, card := range cards {
		actions = append(actions, Action{card.ID, card.Name})
	}

	json.NewEncoder(w).Encode(map[string][]Action{"data": actions})
}

func main() {
	http.HandleFunc("/actions", actions)

	fmt.Printf("Listening on port 8080\n")
	http.ListenAndServe(":8080", nil)
}
