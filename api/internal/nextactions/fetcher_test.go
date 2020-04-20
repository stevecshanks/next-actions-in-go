package nextactions

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/stevecshanks/next-actions-in-go/api/internal/config"
	"github.com/stevecshanks/next-actions-in-go/api/internal/nextactions/mock_nextactions"
	"github.com/stevecshanks/next-actions-in-go/api/internal/trello"
)

func testConfig() *config.Config {
	return &config.Config{
		TrelloNextActionsListID: "nextActionsListId",
		TrelloProjectsListID:    "projectsListId",
	}
}

func testImageURL(size string) url.URL {
	imageURL, _ := url.Parse(fmt.Sprintf("https://trello-backgrounds.s3.amazonaws.com/SharedBackground/%s.jpg", size))
	return *imageURL
}

func testBoard() *trello.Board {
	return &trello.Board{
		ID:   "boardId",
		Name: "My Project",
		Preferences: trello.Preferences{
			BackgroundImages: []trello.BackgroundImage{
				{URL: testImageURL("75x100")},
				{URL: testImageURL("144x192")},
			},
		},
	}
}

func TestOwnedCardsAreReturnedAsActions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_nextactions.NewMockTrelloClient(ctrl)

	cardURL, _ := url.Parse("https://example.com")
	ownedCards := []trello.Card{{ID: "an id", Name: "a name", URL: *cardURL, BoardID: "boardId"}}
	mockClient.EXPECT().OwnedCards().Return(ownedCards, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Any()).Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().GetBoard("boardId").Return(testBoard(), nil).AnyTimes()

	fetcher := Fetcher{mockClient, testConfig()}
	actions, err := fetcher.Fetch()

	expectedActions := []Action{
		{ID: "an id", Name: "a name", URL: *cardURL, ImageURL: testImageURL("75x100"), ProjectName: "My Project"},
	}

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(expectedActions) != len(actions) {
		t.Errorf("Unexpected number of actions returned, expected %d and got %d", len(expectedActions), len(actions))
	}
	if expectedActions[0] != actions[0] {
		t.Errorf("Incorrect actions returned, expected %+v but got %+v", expectedActions, actions)
	}
}

func TestErrorWithOwnedCardsReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_nextactions.NewMockTrelloClient(ctrl)

	expectedError := fmt.Errorf("an error")
	mockClient.EXPECT().OwnedCards().Return(nil, expectedError).AnyTimes()

	fetcher := Fetcher{mockClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err)
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}

func TestCardsInNextActionsListAreReturnedAsActions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_nextactions.NewMockTrelloClient(ctrl)

	nextActionsCards := []trello.Card{{ID: "an id", Name: "a name", BoardID: "boardId"}}
	mockClient.EXPECT().OwnedCards().Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("nextActionsListId")).Return(nextActionsCards, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("projectsListId")).Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().GetBoard("boardId").Return(testBoard(), nil).AnyTimes()

	fetcher := Fetcher{mockClient, testConfig()}
	actions, err := fetcher.Fetch()

	expectedActions := []Action{
		{ID: "an id", Name: "a name", ImageURL: testImageURL("75x100"), ProjectName: "My Project"},
	}

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(expectedActions) != len(actions) {
		t.Errorf("Unexpected number of actions returned, expected %d and got %d", len(expectedActions), len(actions))
	}
	if expectedActions[0] != actions[0] {
		t.Errorf("Incorrect actions returned, expected %+v but got %+v", expectedActions, actions)
	}
}

func TestErrorWithCardsOnListReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_nextactions.NewMockTrelloClient(ctrl)

	expectedError := fmt.Errorf("an error")
	mockClient.EXPECT().OwnedCards().Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("nextActionsListId")).Return(nil, expectedError).AnyTimes()

	fetcher := Fetcher{mockClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err)
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}

func TestErrorWithCardsOnProjectsListReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_nextactions.NewMockTrelloClient(ctrl)

	expectedError := fmt.Errorf("an error")
	mockClient.EXPECT().OwnedCards().Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("nextActionsListId")).Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("projectsListId")).Return(nil, expectedError).AnyTimes()

	fetcher := Fetcher{mockClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err)
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}

func TestMissingDescriptionOnProjectCardReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_nextactions.NewMockTrelloClient(ctrl)

	projectCard := trello.Card{ID: "an id", Name: "a name", Description: "invalid"}
	mockClient.EXPECT().OwnedCards().Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("nextActionsListId")).Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("projectsListId")).Return([]trello.Card{projectCard}, nil).AnyTimes()

	fetcher := Fetcher{mockClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}

func TestErrorWithListsOnBoardReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_nextactions.NewMockTrelloClient(ctrl)

	projectCard := trello.Card{ID: "an id", Name: "a name", Description: "https://trello.com/b/broken/a-broken-card"}
	expectedError := fmt.Errorf("an error")
	mockClient.EXPECT().OwnedCards().Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("nextActionsListId")).Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("projectsListId")).Return([]trello.Card{projectCard}, nil).AnyTimes()
	mockClient.EXPECT().ListsOnBoard(gomock.Eq("broken")).Return(nil, expectedError).AnyTimes()

	fetcher := Fetcher{mockClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err)
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}

func TestMissingTodoListOnProjectBoardReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_nextactions.NewMockTrelloClient(ctrl)

	projectCard := trello.Card{ID: "an id", Name: "a name", Description: "https://trello.com/b/empty"}
	mockClient.EXPECT().OwnedCards().Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("nextActionsListId")).Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("projectsListId")).Return([]trello.Card{projectCard}, nil).AnyTimes()
	mockClient.EXPECT().ListsOnBoard(gomock.Eq("empty")).Return([]trello.List{}, nil).AnyTimes()

	fetcher := Fetcher{mockClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}

func TestErrorWithTodoListReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_nextactions.NewMockTrelloClient(ctrl)

	projectCard := trello.Card{ID: "an id", Name: "a name", Description: "https://trello.com/b/aBoardId"}
	todoList := trello.List{ID: "todoListId", Name: "Todo"}
	anotherList := trello.List{ID: "anotherListId", Name: "Another"}
	expectedError := fmt.Errorf("an error")
	mockClient.EXPECT().OwnedCards().Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("nextActionsListId")).Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("projectsListId")).Return([]trello.Card{projectCard}, nil).AnyTimes()
	mockClient.EXPECT().ListsOnBoard(gomock.Eq("aBoardId")).Return([]trello.List{todoList, anotherList}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("todoListId")).Return(nil, expectedError).AnyTimes()

	fetcher := Fetcher{mockClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err)
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}

func TestEmptyTodoListDoesNotReturnAnAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_nextactions.NewMockTrelloClient(ctrl)

	projectCard := trello.Card{ID: "an id", Name: "a name", Description: "https://trello.com/b/aBoardId"}
	todoList := trello.List{ID: "todoListId", Name: "Todo"}
	mockClient.EXPECT().OwnedCards().Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("nextActionsListId")).Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("projectsListId")).Return([]trello.Card{projectCard}, nil).AnyTimes()
	mockClient.EXPECT().ListsOnBoard(gomock.Eq("aBoardId")).Return([]trello.List{todoList}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("todoListId")).Return([]trello.Card{}, nil).AnyTimes()

	fetcher := Fetcher{mockClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(actions) != 0 {
		t.Errorf("Unexpected number of actions returned, expected %d and got %d", 0, len(actions))
	}
}

func TestFirstTodoListItemsAreReturnedAsActions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_nextactions.NewMockTrelloClient(ctrl)

	projectCard := trello.Card{ID: "an id", Name: "a name", Description: "https://trello.com/b/aBoardId"}
	todoList := trello.List{ID: "todoListId", Name: "Todo"}
	todoListCards := []trello.Card{
		{ID: "an id", Name: "a name", BoardID: "boardId"},
		{ID: "another id", Name: "another name", BoardID: "boardId"},
	}
	mockClient.EXPECT().OwnedCards().Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("nextActionsListId")).Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("projectsListId")).Return([]trello.Card{projectCard}, nil).AnyTimes()
	mockClient.EXPECT().ListsOnBoard(gomock.Eq("aBoardId")).Return([]trello.List{todoList}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("todoListId")).Return(todoListCards, nil).AnyTimes()
	mockClient.EXPECT().GetBoard("boardId").Return(testBoard(), nil).AnyTimes()

	fetcher := Fetcher{mockClient, testConfig()}
	actions, err := fetcher.Fetch()

	expectedActions := []Action{
		{ID: "an id", Name: "a name", ImageURL: testImageURL("75x100"), ProjectName: "My Project"},
	}

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(expectedActions) != len(actions) {
		t.Errorf("Unexpected number of actions returned, expected %d and got %d", len(expectedActions), len(actions))
	}
	if expectedActions[0] != actions[0] {
		t.Errorf("Incorrect actions returned, expected %+v but got %+v", expectedActions, actions)
	}
}

func TestCardDueByDateIsAddedToActions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_nextactions.NewMockTrelloClient(ctrl)

	dueBy, _ := time.Parse(time.RFC3339, "2020-02-12T16:24:00.000Z")
	ownedCards := []trello.Card{{ID: "an id", Name: "a name", DueBy: &dueBy, BoardID: "boardId"}}
	mockClient.EXPECT().OwnedCards().Return(ownedCards, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Any()).Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().GetBoard("boardId").Return(testBoard(), nil).AnyTimes()

	fetcher := Fetcher{mockClient, testConfig()}
	actions, err := fetcher.Fetch()

	expectedActions := []Action{
		{ID: "an id", Name: "a name", DueBy: &dueBy, ImageURL: testImageURL("75x100"), ProjectName: "My Project"},
	}

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(expectedActions) != len(actions) {
		t.Errorf("Unexpected number of actions returned, expected %d and got %d", len(expectedActions), len(actions))
	}
	if expectedActions[0] != actions[0] {
		t.Errorf("Incorrect actions returned, expected %+v but got %+v", expectedActions, actions)
	}
}
