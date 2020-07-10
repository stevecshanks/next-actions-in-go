// Package trello provides a client to interact with the Trello API
package trello

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

// APIBaseURL is the base URL for the Trello API
const APIBaseURL = "https://api.trello.com/1"

// BoardBaseURL is the base URL for Trello boards
const BoardBaseURL = "https://trello.com/b/"

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

// BoardPath returns the path on the Trello API server where a board can be queried
func BoardPath(boardID string) string {
	return fmt.Sprintf("/boards/%s", boardID)
}

// Client is used to interact with the Trello API
type Client struct {
	Key   string
	Token string
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

// GetBoard will return the board with the specified ID
func (c *Client) GetBoard(boardID string) (*Board, error) {
	return c.getBoard(BoardPath(boardID))
}

func (c *Client) getCards(relativePath string) ([]Card, error) {
	resp, err := c.get(relativePath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
	defer resp.Body.Close()

	lists := make([]List, 0)
	if err := json.NewDecoder(resp.Body).Decode(&lists); err != nil {
		return nil, err
	}

	return lists, nil
}

func (c *Client) getBoard(relativePath string) (*Board, error) {
	resp, err := c.get(relativePath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	board := Board{}
	if err := json.NewDecoder(resp.Body).Decode(&board); err != nil {
		return nil, err
	}

	return &board, nil
}

func (c *Client) get(relativePath string) (*http.Response, error) {
	fmt.Printf("Making request to %s\n", relativePath)
	client := &http.Client{}

	relativeURL, err := url.Parse(relativePath)
	if err != nil {
		return nil, err
	}

	queryParameters := relativeURL.Query()
	queryParameters.Add("key", c.Key)
	queryParameters.Add("token", c.Token)

	baseURL, _ := url.Parse(APIBaseURL)
	fullURL := baseURL.ResolveReference(&url.URL{
		Path:     path.Join(baseURL.Path, relativeURL.Path),
		RawQuery: queryParameters.Encode(),
	})

	req, err := http.NewRequest("GET", fullURL.String(), nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("request to %s returned status code %d", relativePath, response.StatusCode)
	}

	return response, nil
}
