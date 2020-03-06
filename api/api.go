package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stevecshanks/next-actions-in-go/api/config"
	"github.com/stevecshanks/next-actions-in-go/api/trello"
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
	config, err := config.FromEnvironment()
	if err != nil {
		handleError(w, err)
	}
	client := trello.Client{
		BaseURL: config.TrelloBaseURL,
		Key:     config.TrelloKey,
		Token:   config.TrelloToken,
	}

	ownedCards, err := client.ListOwnedCards()
	if err != nil {
		handleError(w, err)
	}
	nextActionListCards, err := client.ListCardsOnList(config.TrelloNextActionsListID)
	if err != nil {
		handleError(w, err)
	}

	client.ListCardsOnList(config.TrelloProjectsListID)

	actions := make([]Action, 0)
	for _, card := range append(ownedCards, nextActionListCards...) {
		actions = append(actions, Action{card.ID, card.Name})
	}

	json.NewEncoder(w).Encode(map[string][]Action{"data": actions})
}

func main() {
	http.HandleFunc("/actions", actions)

	fmt.Printf("Listening on port 8080\n")
	http.ListenAndServe(":8080", nil)
}
