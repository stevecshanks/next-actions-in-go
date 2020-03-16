package main

import "github.com/stevecshanks/next-actions-in-go/api/trello"

type TrelloClient interface {
	OwnedCards() ([]trello.Card, error)
}

type Fetcher struct {
	client TrelloClient
}

func (f *Fetcher) Fetch() ([]Action, error) {
	ownedCards, err := f.client.OwnedCards()
	if err != nil {
		return nil, err
	}

	allCards := ownedCards

	actions := make([]Action, 0)
	for _, card := range allCards {
		actions = append(actions, Action{card.ID, card.Name})
	}

	return actions, nil
}
