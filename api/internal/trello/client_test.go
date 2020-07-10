package trello

import (
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestClientOwnedCardsReturnsExpectedResponse(t *testing.T) {
	mockServer := CreateMockServer("some key", "some token")
	defer TeardownMockServer()

	mockServer.AddFileResponse(OwnedCardsPath(), "./testdata/my_cards_response.json")

	client := Client{"some key", "some token"}

	cards, err := client.OwnedCards()
	if err != nil {
		t.Errorf("OwnedCards returned error: %s", err)
	}

	expectedDueBy, _ := time.Parse(time.RFC3339, "2020-01-01T10:30:00.000Z")
	expectedURL1, _ := url.Parse("https://trello.com/c/abcd1234/10-my-first-card")
	expectedCard1 := Card{
		ID:          "myFirstCardId",
		Name:        "My First Action",
		Description: "",
		DueBy:       &expectedDueBy,
		URL:         *expectedURL1,
		BoardID:     "myBoardId",
	}
	expectedURL2, _ := url.Parse("https://trello.com/c/bcde2345/11-my-second-card")
	expectedCard2 := Card{
		ID:          "mySecondCardId",
		Name:        "My Second Action",
		Description: "",
		URL:         *expectedURL2,
		BoardID:     "myBoardId",
	}

	assertCardsMatchExpected(t, cards, []Card{expectedCard1, expectedCard2})
}

func TestClientCardsOnList(t *testing.T) {
	mockServer := CreateMockServer("some key", "some token")
	defer TeardownMockServer()

	mockServer.AddFileResponse(CardsOnListPath("123"), "./testdata/next_actions_list_response.json")

	client := Client{"some key", "some token"}

	cards, err := client.CardsOnList("123")
	if err != nil {
		t.Errorf("CardsOnList returned error: %s", err)
	}

	expectedDueBy, _ := time.Parse(time.RFC3339, "2020-01-15T10:29:59.000Z")
	expectedURL, _ := url.Parse("https://trello.com/c/cdef3456/33-my-third-card")
	expectedCard1 := Card{
		ID:          "todoCardId",
		Name:        "Todo Action",
		Description: "a description",
		DueBy:       &expectedDueBy,
		URL:         *expectedURL,
		BoardID:     "myBoardId",
	}

	assertCardsMatchExpected(t, cards, []Card{expectedCard1})
}

func TestClientListsOnBoard(t *testing.T) {
	mockServer := CreateMockServer("some key", "some token")
	defer TeardownMockServer()

	mockServer.AddFileResponse(ListsOnBoardPath("789"), "./testdata/board_lists_response.json")

	client := Client{"some key", "some token"}

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
	mockServer := CreateMockServer("some key", "some token")
	defer TeardownMockServer()

	mockServer.AddFileResponse(BoardPath("myBoardId"), "./testdata/board_response.json")

	client := Client{"some key", "some token"}

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
	expectedBoard := Board{"myBoardId", "My Project", Preferences{backgroundImages}}
	if expectedBoard.ID != board.ID || expectedBoard.Name != board.Name {
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

func TestClientHandlesHTTPErrors(t *testing.T) {
	CreateMockServer("some key", "some token")
	defer TeardownMockServer()

	client := Client{"some key", "some token"}

	_, err := client.GetBoard("myBoardId")
	if err == nil {
		t.Error("Client did not return 404 error", err)
	}

	expectedError := fmt.Errorf("request to %s returned status code %d", BoardPath("myBoardId"), 404)
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error %s, got %s", expectedError, err)
	}
}

func cardsAreEqual(card, other *Card) bool {
	return (card.ID == other.ID &&
		card.Name == other.Name &&
		card.Description == other.Description &&
		((card.DueBy == nil && other.DueBy == nil) || card.DueBy.Equal(*other.DueBy)) &&
		card.URL.String() == other.URL.String() &&
		card.BoardID == other.BoardID)
}

func assertCardsMatchExpected(t *testing.T, cards, expectedCards []Card) {
	if len(expectedCards) != len(cards) {
		t.Fatalf("Unexpected number of card returned, expected %d and got %d", len(expectedCards), len(cards))
	}
	for i := range cards {
		if !cardsAreEqual(&expectedCards[i], &cards[i]) {
			t.Errorf("Expected card %d to be %+v but got %+v", i, expectedCards[i], cards[i])
		}
	}
}
