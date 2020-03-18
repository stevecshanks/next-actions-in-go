package nextactions

import (
	"encoding/json"
	"time"
)

// Action represents a "next action" in GTD
type Action struct {
	ID    string     `json:"id"`
	Name  string     `json:"name"`
	DueBy *time.Time `json:"dueBy"`
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
