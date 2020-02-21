package trello

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
)

// OwnerCardsPath is the path on the Trello API server where a list of owned cards can be queried
const OwnerCardsPath = "/members/me/cards"

// Card represents a Trello card returned via the API
type Card struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Client is used to interact with the Trello API
type Client struct {
	BaseURL *url.URL
	Key     string
	Token   string
}

// ListOwnedCards will return the cards this user is a member of
func (c *Client) ListOwnedCards() ([]Card, error) {
	client := &http.Client{}

	relativeURL := &url.URL{Path: path.Join(c.BaseURL.Path, OwnerCardsPath)}
	fullURL := c.BaseURL.ResolveReference(relativeURL)
	req, err := http.NewRequest("GET", fullURL.String(), nil)
	if err != nil {
		return nil, err
	}

	queryParameters := req.URL.Query()
	queryParameters.Add("key", c.Key)
	queryParameters.Add("token", c.Token)
	req.URL.RawQuery = queryParameters.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	cards := make([]Card, 0)
	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		return nil, err
	}

	return cards, nil
}
