package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

// Card represents a Trello card returned via the API
type Card struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func handleError(w http.ResponseWriter, err error) {
	fmt.Printf("Error: %s\n", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func actions(w http.ResponseWriter, req *http.Request) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/members/me/cards", os.Getenv("TRELLO_BASE_URL")), nil)
	if err != nil {
		handleError(w, err)
	}

	queryParameters := req.URL.Query()
	queryParameters.Add("key", os.Getenv("TRELLO_KEY"))
	queryParameters.Add("token", os.Getenv("TRELLO_TOKEN"))
	req.URL.RawQuery = queryParameters.Encode()

	resp, err := client.Do(req)
	if err != nil {
		handleError(w, err)
	}

	cards := make([]Card, 0)
	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
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
