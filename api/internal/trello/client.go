package trello

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
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

// BoardPath returns the path on the Trello API server where a board can be queried
func BoardPath(boardID string) string {
	return fmt.Sprintf("/boards/%s", boardID)
}

type URL struct {
	url.URL
}

func (u *URL) UnmarshalJSON(data []byte) error {
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

type CardAlias Card

func (c *Card) UnmarshalJSON(data []byte) error {
	var jsonCard JSONCard
	if err := json.Unmarshal(data, &jsonCard); err != nil {
		return err
	}
	*c = jsonCard.Card()
	return nil
}

type JSONCard struct {
	CardAlias
	URL URL `json:"url"`
}

func (jc *JSONCard) Card() Card {
	card := Card(jc.CardAlias)
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

func (b *BackgroundImage) UnmarshalJSON(data []byte) error {
	var jsonBackgroundImage JSONBackgroundImage
	if err := json.Unmarshal(data, &jsonBackgroundImage); err != nil {
		return err
	}
	*b = BackgroundImage{jsonBackgroundImage.URL.URL}
	return nil
}

type JSONBackgroundImage struct {
	URL URL `json:"url"`
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
