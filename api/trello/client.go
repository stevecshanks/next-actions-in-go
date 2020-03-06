package trello

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

// OwnedCardsPath returns the path on the Trello API server where a list of owned cards can be queried
func OwnedCardsPath() string {
	return "/members/me/cards"
}

// CardsOnListPath returns the path on the Trello API server where cards on a list can be queried
func CardsOnListPath(listID string) string {
	return fmt.Sprintf("/lists/%s/cards", listID)
}

// ListsOnBoardPath returns the path on the Trello API server where lists on a board can be queried
func ListsOnBoardPath(boardID string) string {
	return fmt.Sprintf("/boards/%s/lists", boardID)
}

// Card represents a Trello card returned via the API
type Card struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
}

// List represents a Trello list returned via the API
type List struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Client is used to interact with the Trello API
type Client struct {
	BaseURL *url.URL
	Key     string
	Token   string
}

// OwnedCards will return the cards this user is a member of
func (c *Client) OwnedCards() ([]Card, error) {
	return c.getCards(OwnedCardsPath())
}

// CardsOnList will return the cards on the specified list
func (c *Client) CardsOnList(listID string) ([]Card, error) {
	return c.getCards(CardsOnListPath(listID))
}

// ListsOnBoard will return the lists on the specified board
func (c *Client) ListsOnBoard(boardID string) ([]List, error) {
	return c.getLists(ListsOnBoardPath(boardID))
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

func (c *Client) getLists(relativePath string) ([]List, error) {
	resp, err := c.get(relativePath)
	if err != nil {
		return nil, err
	}

	lists := make([]List, 0)
	if err := json.NewDecoder(resp.Body).Decode(&lists); err != nil {
		return nil, err
	}

	return lists, nil
}

func (c *Client) get(relativePath string) (*http.Response, error) {
	client := &http.Client{}

	relativeURL, err := url.Parse(relativePath)
	if err != nil {
		return nil, err
	}

	queryParameters := relativeURL.Query()
	queryParameters.Add("key", c.Key)
	queryParameters.Add("token", c.Token)

	fullURL := c.BaseURL.ResolveReference(&url.URL{
		Path:     path.Join(c.BaseURL.Path, relativeURL.Path),
		RawQuery: queryParameters.Encode(),
	})

	req, err := http.NewRequest("GET", fullURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}
