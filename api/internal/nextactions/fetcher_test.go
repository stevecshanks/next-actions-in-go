package nextactions

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/stevecshanks/next-actions-in-go/api/internal/config"
	"github.com/stevecshanks/next-actions-in-go/api/internal/nextactions/mock_nextactions"
	"github.com/stevecshanks/next-actions-in-go/api/internal/trello"
)

func testConfig() *config.Config {
	return &config.Config{
		TrelloNextActionsListID: "nextActionsListId",
	}
}

func TestOwnedCardsAreReturnedAsActions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_nextactions.NewMockTrelloClient(ctrl)

	ownedCards := []trello.Card{{ID: "an id", Name: "a name", Description: ""}}
	mockClient.EXPECT().OwnedCards().Return(ownedCards, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Any()).Return([]trello.Card{}, nil).AnyTimes()

	fetcher := Fetcher{mockClient, testConfig()}
	actions, err := fetcher.Fetch()

	expectedActions := []Action{
		{"an id", "a name"},
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

	nextActionsCards := []trello.Card{{ID: "an id", Name: "a name", Description: ""}}
	mockClient.EXPECT().OwnedCards().Return([]trello.Card{}, nil).AnyTimes()
	mockClient.EXPECT().CardsOnList(gomock.Eq("nextActionsListId")).Return(nextActionsCards, nil).AnyTimes()

	fetcher := Fetcher{mockClient, testConfig()}
	actions, err := fetcher.Fetch()

	expectedActions := []Action{
		{"an id", "a name"},
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
