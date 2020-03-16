package main

import (
	"fmt"
	"testing"

	"github.com/stevecshanks/next-actions-in-go/api/trello"
)

type mockClient struct {
	OwnedCardsReturn []trello.Card
	OwnedCardsError  error
}

func (m *mockClient) OwnedCards() ([]trello.Card, error) {
	return m.OwnedCardsReturn, m.OwnedCardsError
}

func TestOwnedCardsAreReturnedAsActions(t *testing.T) {
	expectedActions := []Action{
		{"an id", "a name"},
	}

	fetcher := Fetcher{&mockClient{
		OwnedCardsReturn: []trello.Card{{ID: "an id", Name: "a name", Description: ""}},
	}}

	actions, err := fetcher.Fetch()

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
	expectedError := fmt.Errorf("an error")

	fetcher := Fetcher{&mockClient{
		OwnedCardsError: expectedError,
	}}

	actions, err := fetcher.Fetch()

	if err != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err)
	}
	if actions != nil {
		t.Errorf("Expected no actions, got %+v", actions)
	}
}
