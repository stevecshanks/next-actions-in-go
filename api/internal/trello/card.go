package trello // nolint:golint // package comment is in another file

import (
	"encoding/json"
	"net/url"
	"time"
)

type urlWrapper struct {
	url.URL
}

func (u *urlWrapper) UnmarshalJSON(data []byte) error {
	var jsonURL string
	if err := json.Unmarshal(data, &jsonURL); err != nil {
		return err
	}

	parsedURL, err := url.Parse(jsonURL)
	if err != nil {
		return err
	}

	u.URL = *parsedURL

	return nil
}

// Card represents a Trello card returned via the API
type Card struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	DueBy   *time.Time `json:"due"`
	URL     url.URL    `json:"-"`
	BoardID string     `json:"idBoard"`
}

type cardAlias Card

// UnmarshalJSON converts JSON data into a Card
func (c *Card) UnmarshalJSON(data []byte) error {
	var jsonCard jsonCard
	if err := json.Unmarshal(data, &jsonCard); err != nil {
		return err
	}
	*c = jsonCard.Card()
	return nil
}

type jsonCard struct {
	cardAlias
	URL urlWrapper `json:"url"`
}

func (jc *jsonCard) Card() Card {
	card := Card(jc.cardAlias)
	card.URL = jc.URL.URL
	return card
}
