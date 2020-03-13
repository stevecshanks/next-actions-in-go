package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

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
	cfg, err := config.FromEnvironment()
	if err != nil {
		handleError(w, err)
	}
	client := trello.Client{
		BaseURL: cfg.TrelloBaseURL,
		Key:     cfg.TrelloKey,
		Token:   cfg.TrelloToken,
	}

	ownedCards, err := client.OwnedCards()
	if err != nil {
		handleError(w, err)
	}
	nextActionListCards, err := client.CardsOnList(cfg.TrelloNextActionsListID)
	if err != nil {
		handleError(w, err)
	}

	projectCards, err := client.CardsOnList(cfg.TrelloProjectsListID)
	if err != nil {
		handleError(w, err)
	}

	// Note: This is NOT the same as the API base URL
	boardsURLPath := "https://trello.com/b/"
	boardIDRegex, err := regexp.Compile(regexp.QuoteMeta(boardsURLPath) + `(\w+).*`)
	if err != nil {
		handleError(w, err)
	}

	allCards := append(ownedCards, nextActionListCards...)
	for _, projectCard := range projectCards {
		projectBoardID := boardIDRegex.FindStringSubmatch(projectCard.Description)[1]
		projectLists, err := client.ListsOnBoard(projectBoardID)
		if err != nil {
			handleError(w, err)
		}

		for _, list := range projectLists {
			if list.Name == "Todo" {
				projectTodoListCards, err := client.CardsOnList(list.ID)
				if err != nil {
					handleError(w, err)
				}
				if len(projectTodoListCards) > 0 {
					allCards = append(allCards, projectTodoListCards[0])
				}
				break
			}
		}
	}

	actions := make([]Action, 0)
	for _, card := range allCards {
		actions = append(actions, Action{card.ID, card.Name})
	}

	if err := json.NewEncoder(w).Encode(map[string][]Action{"data": actions}); err != nil {
		handleError(w, err)
	}
}

func main() {
	http.HandleFunc("/actions", actions)

	fmt.Printf("Listening on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
