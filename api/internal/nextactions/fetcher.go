package nextactions

import "github.com/stevecshanks/next-actions-in-go/api/internal/trello"

type TrelloClient interface {
	OwnedCards() ([]trello.Card, error)
}

type Fetcher struct {
	Client TrelloClient
}

func (f *Fetcher) Fetch() ([]Action, error) {
	ownedCards, err := f.Client.OwnedCards()
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
