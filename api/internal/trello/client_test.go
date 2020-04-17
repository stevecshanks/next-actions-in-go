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
	expectedURL1, _ := url.Parse("https://trello.com/c/abcd1234/10-my-first-card")
	expectedCard1 := Card{
		ID:          "myFirstCardId",
		Name:        "My First Card",
		Description: "",
		DueBy:       &expectedDueBy,
		URL:         *expectedURL1,
	}
	if !cardsAreEqual(&cards[0], &expectedCard1) {
		t.Errorf(fmt.Sprintf("OwnedCards returned incorrect card, expected %+v got %+v", expectedCard1, cards[0]))
	}
	expectedURL2, _ := url.Parse("https://trello.com/c/bcde2345/11-my-second-card")
	expectedCard2 := Card{
		ID:          "mySecondCardId",
		Name:        "My Second Card",
		Description: "",
		URL:         *expectedURL2,
	}
	if !cardsAreEqual(&cards[1], &expectedCard2) {
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
	expectedURL, _ := url.Parse("https://trello.com/c/cdef3456/33-my-third-card")
	expectedCard1 := Card{
		ID:          "todoCardId",
		Name:        "Todo Card",
		Description: "a description",
		URL:         *expectedURL,
	}
	if !cardsAreEqual(&cards[0], &expectedCard1) {
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

func TestClientGetBoard(t *testing.T) {
	mockServer := CreateMockServer("https://api.trello.com/1", "some key", "some token")
	defer TeardownMockServer()

	mockServer.AddFileResponse(BoardPath("myBoardId"), "./testdata/board_response.json")

	baseURL, _ := url.Parse("https://api.trello.com/1")
	client := Client{baseURL, "some key", "some token"}

	board, err := client.GetBoard("myBoardId")
	if err != nil {
		t.Errorf("GetBoard returned error: %s", err)
	}

	backgroundURL1, _ := url.Parse("https://trello-backgrounds.s3.amazonaws.com/SharedBackground/75x100.jpg")
	backgroundURL2, _ := url.Parse("https://trello-backgrounds.s3.amazonaws.com/SharedBackground/144x192.jpg")
	backgroundImages := []BackgroundImage{
		{*backgroundURL1},
		{*backgroundURL2},
	}
	expectedBoard := Board{"myBoardId", Preferences{backgroundImages}}
	if expectedBoard.ID != board.ID {
		t.Errorf(fmt.Sprintf("GetBoard returned incorrect board, expected %+v got %+v", expectedBoard, board))
	}
	if len(expectedBoard.Preferences.BackgroundImages) != len(board.Preferences.BackgroundImages) {
		t.Fatalf(fmt.Sprintf("GetBoard returned incorrect board, expected %+v got %+v", expectedBoard, board))
	}
	for i, expectedBackgroundImage := range expectedBoard.Preferences.BackgroundImages {
		if expectedBackgroundImage.URL.String() != board.Preferences.BackgroundImages[i].URL.String() {
			t.Errorf(fmt.Sprintf("GetBoard returned incorrect board, expected %+v got %+v", expectedBoard, board))
		}
	}
}

func cardsAreEqual(card, other *Card) bool {
	return (card.ID == other.ID &&
		card.Name == other.Name &&
		card.Description == other.Description &&
		((card.DueBy == nil && other.DueBy == nil) || card.DueBy.Equal(*other.DueBy)) &&
		card.URL.String() == other.URL.String())
}
