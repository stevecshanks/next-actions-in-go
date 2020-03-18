package trello

import (
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestClientOwnedCardsReturnsExpectedResponse(t *testing.T) {
	mockServer := CreateMockServer("https://api.trello.com/1", "some key", "some token")
	defer TeardownMockServer()

	mockServer.AddFileResponse(OwnedCardsPath(), "./testdata/my_cards_response.json")

	baseURL, _ := url.Parse("https://api.trello.com/1")
	client := Client{baseURL, "some key", "some token"}

	cards, err := client.OwnedCards()
	if err != nil {
		t.Errorf("OwnedCards returned error: %s", err)
	}
	if len(cards) != 2 {
		t.Fatalf("OwnedCards returned %d cards, expected %d", len(cards), 2)
	}
	expectedDueBy, _ := time.Parse(time.RFC3339, "2020-02-12T16:24:00.000Z")
	expectedCard1 := Card{
		ID:          "myFirstCardId",
		Name:        "My First Card",
		Description: "",
		DueBy:       &expectedDueBy,
	}
	if cards[0].ID != expectedCard1.ID ||
		cards[0].Name != expectedCard1.Name ||
		cards[0].Description != expectedCard1.Description ||
		!cards[0].DueBy.Equal(*expectedCard1.DueBy) {
		t.Errorf(fmt.Sprintf("OwnedCards returned incorrect card, expected %+v got %+v", expectedCard1, cards[0]))
	}
	expectedCard2 := Card{
		ID:          "mySecondCardId",
		Name:        "My Second Card",
		Description: "",
	}
	if cards[1] != expectedCard2 {
		t.Errorf(fmt.Sprintf("OwnedCards returned incorrect card, expected %+v got %+v", expectedCard2, cards[1]))
	}
}

func TestClientCardsOnList(t *testing.T) {
	mockServer := CreateMockServer("https://api.trello.com/1", "some key", "some token")
	defer TeardownMockServer()

	mockServer.AddFileResponse(CardsOnListPath("123"), "./testdata/next_actions_list_response.json")

	baseURL, _ := url.Parse("https://api.trello.com/1")
	client := Client{baseURL, "some key", "some token"}

	cards, err := client.CardsOnList("123")
	if err != nil {
		t.Errorf("CardsOnList returned error: %s", err)
	}
	if len(cards) != 1 {
		t.Fatalf("CardsOnList returned %d cards, expected %d", len(cards), 2)
	}
	expectedCard1 := Card{
		ID:          "todoCardId",
		Name:        "Todo Card",
		Description: "a description",
	}
	if cards[0] != expectedCard1 {
		t.Errorf(fmt.Sprintf("CardsOnList returned incorrect card, expected %+v got %+v", expectedCard1, cards[0]))
	}
}

func TestClientListsOnBoard(t *testing.T) {
	mockServer := CreateMockServer("https://api.trello.com/1", "some key", "some token")
	defer TeardownMockServer()

	mockServer.AddFileResponse(ListsOnBoardPath("789"), "./testdata/board_lists_response.json")

	baseURL, _ := url.Parse("https://api.trello.com/1")
	client := Client{baseURL, "some key", "some token"}

	lists, err := client.ListsOnBoard("789")
	if err != nil {
		t.Errorf("ListsOnBoard returned error: %s", err)
	}
	if len(lists) != 2 {
		t.Fatalf("ListsOnBoard returned %d lists, expected %d", len(lists), 2)
	}
	expectedList1 := List{"inboxListId", "Inbox"}
	expectedList2 := List{"todoListId", "Todo"}
	if lists[0] != expectedList1 {
		t.Errorf(fmt.Sprintf("ListsOnBoard returned incorrect list, expected %+v got %+v", expectedList1, lists[0]))
	}
	if lists[1] != expectedList2 {
		t.Errorf(fmt.Sprintf("ListsOnBoard returned incorrect list, expected %+v got %+v", expectedList2, lists[1]))
	}
}
