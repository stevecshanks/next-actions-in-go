package nextactions // nolint:golint // package comment is in another file

import (
	"encoding/json"
	"net/url"
	"time"
)

// Action represents a "next action" in GTD
type Action struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	DueBy       *time.Time `json:"dueBy"`
	URL         url.URL    `json:"-"`
	ImageURL    *url.URL   `json:"-"`
	ProjectName string     `json:"projectName"`
}

type actionAlias Action

// MarshalJSON returns a JSON representation of an Action
func (a *Action) MarshalJSON() ([]byte, error) {
	var imageURL *string = nil
	if a.ImageURL != nil {
		var u = a.ImageURL.String()
		imageURL = &u
	}
	return json.Marshal(jsonAction{
		actionAlias: actionAlias(*a),
		Type:        "actions",
		URL:         a.URL.String(),
		ImageURL:    imageURL,
	})
}

type jsonAction struct {
	Type string `json:"type"` // Required by JSON-API
	actionAlias
	URL      string  `json:"url"`
	ImageURL *string `json:"imageUrl"`
}
