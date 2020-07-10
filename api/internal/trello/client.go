// Package trello provides a client to interact with the Trello API
package trello

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
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
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"desc"`
	DueBy       *time.Time `json:"due"`
	URL         url.URL    `json:"-"`
	BoardID     string     `json:"idBoard"`
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

// List represents a Trello list returned via the API
type List struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

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
