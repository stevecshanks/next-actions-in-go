package trello

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestClientListOwnedCardsReturnsExpectedResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	bytes, err := ioutil.ReadFile("./testdata/my_cards_response.json")
	if err != nil {
		t.Fatal(err)
	}
	httpmock.RegisterResponderWithQuery(
		"GET",
		"https://api.trello.com/1/members/me/cards",
		"key=some+key&token=some+token",
		httpmock.NewBytesResponder(200, bytes),
	)

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
