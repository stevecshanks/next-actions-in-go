package trello

import (
	"fmt"
	"net/url"
	"testing"
)

func TestClientListOwnedCardsReturnsExpectedResponse(t *testing.T) {
	mockServer := CreateMockServer("https://api.trello.com/1", "some key", "some token")
	defer TeardownMockServer()

	mockServer.AddFileResponse("members/me/cards", "./testdata/my_cards_response.json")

	baseURL, _ := url.Parse("https://api.trello.com/1")
	client := Client{baseURL, "some key", "some token"}

	cards, err := client.ListOwnedCards()
	if err != nil {
		t.Errorf("ListOwnedCards returned error: %s", err)
	}
	if len(cards) != 2 {
		t.Errorf("ListOwnedCards returned %d cards, expected %d", len(cards), 2)
	}
	expectedCard1 := Card{"111111111111111111111111", "My First Card"}
	if cards[0] != expectedCard1 {
		t.Errorf(fmt.Sprintf("ListOwnedCards returned incorrect card, expected %+v got %+v", expectedCard1, cards[0]))
	}
	expectedCard2 := Card{"222222222222222222222222", "My Second Card"}
	if cards[1] != expectedCard2 {
		t.Errorf(fmt.Sprintf("ListOwnedCards returned incorrect card, expected %+v got %+v", expectedCard2, cards[1]))
	}
}

func TestClientListCardsOnList(t *testing.T) {
	mockServer := CreateMockServer("https://api.trello.com/1", "some key", "some token")
	defer TeardownMockServer()

	mockServer.AddFileResponse("lists/123/cards", "./testdata/next_actions_list_response.json")

	baseURL, _ := url.Parse("https://api.trello.com/1")
	client := Client{baseURL, "some key", "some token"}

	cards, err := client.ListCardsOnList("123")
	if err != nil {
		t.Errorf("ListCardsOnList returned error: %s", err)
	}
	if len(cards) != 1 {
		t.Errorf("ListCardsOnList returned %d cards, expected %d", len(cards), 2)
	}
	expectedCard1 := Card{"333333333333333333333333", "My Third Card"}
	if cards[0] != expectedCard1 {
		t.Errorf(fmt.Sprintf("ListCardsOnList returned incorrect card, expected %+v got %+v", expectedCard1, cards[0]))
	}
}
