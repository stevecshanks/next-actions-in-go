package trello

import (
	"fmt"
	"net/url"
	"testing"
)

func TestClientListOwnedCardsReturnsExpectedResponse(t *testing.T) {
	trelloKey := "my_trello_key"
	trelloToken := "my_trello_token"

	server, teardown := CreateMockServer(t, trelloKey, trelloToken)
	defer teardown()

	baseURL, _ := url.Parse(server.URL)
	client := Client{baseURL, trelloKey, trelloToken}

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
