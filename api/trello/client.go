package trello

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

// OwnedCardsPath is the path on the Trello API server where a list of owned cards can be queried
const OwnedCardsPath = "/members/me/cards"

// CardsOnListPathTemplate is a template for the path on the Trello API server where cards on a list can be queried
// Needs to be formatted with the List ID
const CardsOnListPathTemplate = "/lists/%s/cards"

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
	return c.getCards(OwnedCardsPath)
}

// ListCardsOnList will return the cards on the specified list
func (c *Client) ListCardsOnList(listID string) ([]Card, error) {
	return c.getCards(fmt.Sprintf(CardsOnListPathTemplate, listID))
}

func (c *Client) getCards(relativePath string) ([]Card, error) {
	resp, err := c.get(relativePath)
	if err != nil {
		return nil, err
	}

	cards := make([]Card, 0)
	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		return nil, err
	}

	return cards, nil
}

func (c *Client) get(relativePath string) (*http.Response, error) {
	client := &http.Client{}

	relativeURL := &url.URL{Path: path.Join(c.BaseURL.Path, relativePath)}
	fullURL := c.BaseURL.ResolveReference(relativeURL)
	req, err := http.NewRequest("GET", fullURL.String(), nil)
	if err != nil {
		return nil, err
	}

	queryParameters := req.URL.Query()
	queryParameters.Add("key", c.Key)
	queryParameters.Add("token", c.Token)
	req.URL.RawQuery = queryParameters.Encode()

	return client.Do(req)
}
