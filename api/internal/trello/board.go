package trello // nolint:golint // package comment is in another file

import (
	"encoding/json"
	"net/url"
)

// Board represents a Trello board returned via the API
type Board struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Preferences Preferences `json:"prefs"`
}

// Preferences represents the preferences for a Trello board returned via the API
type Preferences struct {
	BackgroundImages []BackgroundImage `json:"backgroundImageScaled"`
}

// BackgroundImage represents the background image of a Trello board returned via the API
type BackgroundImage struct {
	URL url.URL `json:"-"`
}

// UnmarshalJSON converts JSON data into a BackgroundImage
func (b *BackgroundImage) UnmarshalJSON(data []byte) error {
	var jsonBackgroundImage jsonBackgroundImage
	if err := json.Unmarshal(data, &jsonBackgroundImage); err != nil {
		return err
	}
	*b = BackgroundImage{jsonBackgroundImage.URL.URL}
	return nil
}

type jsonBackgroundImage struct {
	URL urlWrapper `json:"url"`
}
