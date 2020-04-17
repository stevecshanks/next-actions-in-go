package nextactions

import (
	"encoding/json"
	"net/url"
	"time"
)

type URL struct {
	url.URL
}

func (u *URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.URL.String())
}

// Action represents a "next action" in GTD
type Action struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	DueBy    *time.Time `json:"dueBy"`
	URL      url.URL    `json:"-"`
	ImageURL url.URL    `json:"-"`
}

type ActionAlias Action

func (a *Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(JSONAction{
		ActionAlias: ActionAlias(*a),
		Type:        "actions",
		URL:         a.URL.String(),
		ImageURL:    a.ImageURL.String(),
	})
}

type JSONAction struct {
	Type string `json:"type"` // Required by JSON-API
	ActionAlias
	URL      string `json:"url"`
	ImageURL string `json:"imageUrl"`
}
