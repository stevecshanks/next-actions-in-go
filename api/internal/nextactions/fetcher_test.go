package nextactions

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stevecshanks/next-actions-in-go/api/internal/config"
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

type fakeTrelloClient struct {
	ownedCards         []trello.Card
	cardsOnLists       map[string][]trello.Card
	listsOnBoards      map[string][]trello.List
	ownedCardsError    error
	cardsOnListErrors  map[string]error
	listsOnBoardErrors map[string]error
}

func (f *fakeTrelloClient) OwnedCards() ([]trello.Card, error) {
	if f.ownedCardsError != nil {
		return nil, f.ownedCardsError
	}
	return f.ownedCards, nil
}

func (f *fakeTrelloClient) CardsOnList(listID string) ([]trello.Card, error) {
	if f.cardsOnListErrors[listID] != nil {
		return nil, f.cardsOnListErrors[listID]
	}
	cards, ok := f.cardsOnLists[listID]
	if !ok {
		return []trello.Card{}, nil
	}
	return cards, nil
}

func (f *fakeTrelloClient) ListsOnBoard(boardID string) ([]trello.List, error) {
	if f.listsOnBoardErrors[boardID] != nil {
		return nil, f.listsOnBoardErrors[boardID]
	}
	lists, ok := f.listsOnBoards[boardID]
	if !ok {
		return []trello.List{}, nil
	}
	return lists, nil
}

func (f *fakeTrelloClient) GetBoard(boardID string) (*trello.Board, error) {
	return &trello.Board{
		ID:   boardID,
		Name: "My Project",
		Preferences: trello.Preferences{
			BackgroundImages: []trello.BackgroundImage{
				{URL: testImageURL("75x100")},
				{URL: testImageURL("144x192")},
			},
		},
	}, nil
}

func (f *fakeTrelloClient) AddOwnedCard(card *trello.Card) {
	f.ownedCards = append(f.ownedCards, *card)
}

func (f *fakeTrelloClient) AddCardOnList(listID string, card *trello.Card) {
	f.cardsOnLists[listID] = append(f.cardsOnLists[listID], *card)
}

func (f *fakeTrelloClient) AddListOnBoard(boardID string, list *trello.List) {
	f.listsOnBoards[boardID] = append(f.listsOnBoards[boardID], *list)
}

func (f *fakeTrelloClient) SetOwnedCardsError(err error) {
	f.ownedCardsError = err
}

func (f *fakeTrelloClient) SetCardsOnListError(listID string, err error) {
	f.cardsOnListErrors[listID] = err
}

func (f *fakeTrelloClient) SetListsOnBoardError(boardID string, err error) {
	f.listsOnBoardErrors[boardID] = err
}

func aFakeTrelloClient() *fakeTrelloClient {
	return &fakeTrelloClient{
		cardsOnLists:       make(map[string][]trello.Card),
		listsOnBoards:      make(map[string][]trello.List),
		cardsOnListErrors:  make(map[string]error),
		listsOnBoardErrors: make(map[string]error),
	}
}

func TestOwnedCardsAreReturnedAsActions(t *testing.T) {
	cardURL, _ := url.Parse("https://example.com")
	ownedCard := trello.Card{ID: "an id", Name: "a name", URL: *cardURL, BoardID: "boardId"}

	fakeClient := aFakeTrelloClient()
	fakeClient.AddOwnedCard(&ownedCard)

	fetcher := Fetcher{fakeClient, testConfig()}
	actions, err := fetcher.Fetch()

	expectedActions := []Action{
		{ID: "an id", Name: "a name", URL: *cardURL, ImageURL: testImageURL("75x100"), ProjectName: "My Project"},
	}

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(expectedActions) != len(actions) {
		t.Fatalf("Unexpected number of actions returned, expected %d and got %d", len(expectedActions), len(actions))
	}
	if expectedActions[0] != actions[0] {
		t.Errorf("Incorrect actions returned, expected %+v but got %+v", expectedActions, actions)
	}
}

func TestErrorWithOwnedCardsReturnsError(t *testing.T) {
	expectedError := fmt.Errorf("an error")

	fakeClient := aFakeTrelloClient()
	fakeClient.SetOwnedCardsError(expectedError)

	fetcher := Fetcher{fakeClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err)
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}

func TestCardsInNextActionsListAreReturnedAsActions(t *testing.T) {
	nextActionsCard := trello.Card{ID: "an id", Name: "a name", BoardID: "boardId"}

	fakeClient := aFakeTrelloClient()
	fakeClient.AddCardOnList("nextActionsListId", &nextActionsCard)

	fetcher := Fetcher{fakeClient, testConfig()}
	actions, err := fetcher.Fetch()

	expectedActions := []Action{
		{ID: "an id", Name: "a name", ImageURL: testImageURL("75x100"), ProjectName: "My Project"},
	}

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(expectedActions) != len(actions) {
		t.Fatalf("Unexpected number of actions returned, expected %d and got %d", len(expectedActions), len(actions))
	}
	if expectedActions[0] != actions[0] {
		t.Errorf("Incorrect actions returned, expected %+v but got %+v", expectedActions, actions)
	}
}

func TestErrorWithCardsOnListReturnsError(t *testing.T) {
	expectedError := fmt.Errorf("an error")

	fakeClient := aFakeTrelloClient()
	fakeClient.SetCardsOnListError("nextActionsListId", expectedError)

	fetcher := Fetcher{fakeClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err)
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}

func TestErrorWithCardsOnProjectsListReturnsError(t *testing.T) {
	expectedError := fmt.Errorf("an error")

	fakeClient := aFakeTrelloClient()
	fakeClient.SetCardsOnListError("projectsListId", expectedError)

	fetcher := Fetcher{fakeClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err)
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}

func TestMissingDescriptionOnProjectCardReturnsError(t *testing.T) {
	brokenProjectCard := trello.Card{ID: "an id", Name: "a name", Description: "invalid"}

	fakeClient := aFakeTrelloClient()
	fakeClient.AddCardOnList("projectsListId", &brokenProjectCard)

	fetcher := Fetcher{fakeClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}

func TestErrorWithListsOnBoardReturnsError(t *testing.T) {
	projectCard := trello.Card{ID: "an id", Name: "a name", Description: "https://trello.com/b/broken/a-broken-card"}
	expectedError := fmt.Errorf("an error")

	fakeClient := aFakeTrelloClient()
	fakeClient.AddCardOnList("projectsListId", &projectCard)
	fakeClient.SetListsOnBoardError("broken", expectedError)

	fetcher := Fetcher{fakeClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err)
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}

func TestMissingTodoListOnProjectBoardReturnsError(t *testing.T) {
	projectCard := trello.Card{ID: "an id", Name: "a name", Description: "https://trello.com/b/empty"}

	fakeClient := aFakeTrelloClient()
	fakeClient.AddCardOnList("projectsListId", &projectCard)

	fetcher := Fetcher{fakeClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}

func TestErrorWithTodoListReturnsError(t *testing.T) {
	projectCard := trello.Card{ID: "an id", Name: "a name", Description: "https://trello.com/b/aBoardId"}
	todoList := trello.List{ID: "todoListId", Name: "Todo"}
	expectedError := fmt.Errorf("an error")

	fakeClient := aFakeTrelloClient()
	fakeClient.AddCardOnList("projectsListId", &projectCard)
	fakeClient.AddListOnBoard("aBoardId", &todoList)
	fakeClient.SetCardsOnListError("todoListId", expectedError)

	fetcher := Fetcher{fakeClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err)
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}

func TestEmptyTodoListDoesNotReturnAnAction(t *testing.T) {
	projectCard := trello.Card{ID: "an id", Name: "a name", Description: "https://trello.com/b/aBoardId"}
	todoList := trello.List{ID: "todoListId", Name: "Todo"}

	fakeClient := aFakeTrelloClient()
	fakeClient.AddCardOnList("projectsListId", &projectCard)
	fakeClient.AddListOnBoard("aBoardId", &todoList)

	fetcher := Fetcher{fakeClient, testConfig()}
	actions, err := fetcher.Fetch()

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(actions) != 0 {
		t.Errorf("Unexpected number of actions returned, expected %d and got %d", 0, len(actions))
	}
}

func TestFirstTodoListItemsAreReturnedAsActions(t *testing.T) {
	projectCard := trello.Card{ID: "an id", Name: "a name", Description: "https://trello.com/b/aBoardId"}
	todoList := trello.List{ID: "todoListId", Name: "Todo"}

	fakeClient := aFakeTrelloClient()
	fakeClient.AddCardOnList("projectsListId", &projectCard)
	fakeClient.AddListOnBoard("aBoardId", &todoList)
	fakeClient.AddCardOnList("todoListId", &trello.Card{ID: "an id", Name: "a name", BoardID: "boardId"})
	fakeClient.AddCardOnList("todoListId", &trello.Card{ID: "another id", Name: "another name", BoardID: "boardId"})

	fetcher := Fetcher{fakeClient, testConfig()}
	actions, err := fetcher.Fetch()

	expectedActions := []Action{
		{ID: "an id", Name: "a name", ImageURL: testImageURL("75x100"), ProjectName: "My Project"},
	}

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(expectedActions) != len(actions) {
		t.Fatalf("Unexpected number of actions returned, expected %d and got %d", len(expectedActions), len(actions))
	}
	if expectedActions[0] != actions[0] {
		t.Errorf("Incorrect actions returned, expected %+v but got %+v", expectedActions, actions)
	}
}

func TestCardDueByDateIsAddedToActions(t *testing.T) {
	dueBy, _ := time.Parse(time.RFC3339, "2020-02-12T16:24:00.000Z")
	ownedCard := trello.Card{ID: "an id", Name: "a name", DueBy: &dueBy, BoardID: "boardId"}

	fakeClient := aFakeTrelloClient()
	fakeClient.AddOwnedCard(&ownedCard)

	fetcher := Fetcher{fakeClient, testConfig()}
	actions, err := fetcher.Fetch()

	expectedActions := []Action{
		{ID: "an id", Name: "a name", DueBy: &dueBy, ImageURL: testImageURL("75x100"), ProjectName: "My Project"},
	}

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(expectedActions) != len(actions) {
		t.Fatalf("Unexpected number of actions returned, expected %d and got %d", len(expectedActions), len(actions))
	}
	if expectedActions[0] != actions[0] {
		t.Errorf("Incorrect actions returned, expected %+v but got %+v", expectedActions, actions)
	}
}
