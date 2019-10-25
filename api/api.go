package main

import (
	"encoding/json"
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
		Action{
			ID:   "action2",
			Name: "Another action",
		},
	}

	json.NewEncoder(w).Encode(map[string][]Action{"data": dummyActions})
}

func main() {
	http.HandleFunc("/actions", actions)

	http.ListenAndServe(":8080", nil)
}
