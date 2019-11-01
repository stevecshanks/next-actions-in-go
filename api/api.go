package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func actions(w http.ResponseWriter, req *http.Request) {
	var dummyActions = []Action{
		Action{
			ID:   "action1",
			Name: "Some action",
		},
	}

	json.NewEncoder(w).Encode(map[string][]Action{"data": dummyActions})
}

func main() {
	http.HandleFunc("/actions", actions)

	fmt.Printf("Listening on port 8080\n")
	http.ListenAndServe(":8080", nil)
}
