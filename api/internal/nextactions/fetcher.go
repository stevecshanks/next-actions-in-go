package nextactions

import (
	"fmt"
	"regexp"

	"github.com/stevecshanks/next-actions-in-go/api/internal/config"
	"github.com/stevecshanks/next-actions-in-go/api/internal/trello"
)

// Note: This is NOT the same as the API base URL
const boardsURLPath = "https://trello.com/b/"

type TrelloClient interface {
	OwnedCards() ([]trello.Card, error)
	CardsOnList(listID string) ([]trello.Card, error)
	ListsOnBoard(boardID string) ([]trello.List, error)
	GetBoard(boardID string) (*trello.Board, error)
}

type Fetcher struct {
	Client TrelloClient
	Config *config.Config
}

func (f *Fetcher) Fetch() ([]Action, error) {
	allCards := make([]trello.Card, 0)

	ownedCards, err := f.fetchOwnedCards()
	if err != nil {
		return nil, err
	}
	allCards = append(allCards, ownedCards...)

	nextActionsCards, err := f.fetchCardsOnNextActionsList()
	if err != nil {
		return nil, err
	}
	allCards = append(allCards, nextActionsCards...)

	projectTodoCards, err := f.fetchProjectTodoListCards()
	if err != nil {
		return nil, err
	}
	allCards = append(allCards, projectTodoCards...)

	boardsByID, err := f.fetchAllBoards(allCards)
	if err != nil {
		return nil, err
	}

	return cardsToActions(allCards, boardsByID), nil
}

func (f *Fetcher) fetchOwnedCards() ([]trello.Card, error) {
	return f.Client.OwnedCards()
}

func (f *Fetcher) fetchCardsOnNextActionsList() ([]trello.Card, error) {
	return f.Client.CardsOnList(f.Config.TrelloNextActionsListID)
}

func (f *Fetcher) fetchProjectTodoListCards() ([]trello.Card, error) {
	allCards := make([]trello.Card, 0)

	projectCards, err := f.Client.CardsOnList(f.Config.TrelloProjectsListID)
	if err != nil {
		return nil, err
	}

	todoListCardsChannel := make(chan []trello.Card)
	errorsChannel := make(chan error)

	for i := range projectCards {
		projectCard := &projectCards[i]
		go f.fetchProjectTodoList(projectCard, todoListCardsChannel, errorsChannel)
	}

	for range projectCards {
		select {
		case todoListCards := <-todoListCardsChannel:
			if len(todoListCards) > 0 {
				allCards = append(allCards, todoListCards[0])
			}
		case err := <-errorsChannel:
			return nil, err
		}
	}

	return allCards, nil
}

func (f *Fetcher) fetchProjectTodoList(
	projectCard *trello.Card,
	todoListCardsChannel chan []trello.Card,
	errorsChannel chan error,
) {
	projectBoardID, err := getProjectBoardID(projectCard)
	if err != nil {
		errorsChannel <- err
		return
	}
	projectLists, err := f.Client.ListsOnBoard(projectBoardID)
	if err != nil {
		errorsChannel <- err
		return
	}

	todoList, err := getTodoList(projectLists)
	if err != nil {
		errorsChannel <- err
		return
	}
	todoListCards, err := f.Client.CardsOnList(todoList.ID)
	if err != nil {
		errorsChannel <- err
		return
	}

	todoListCardsChannel <- todoListCards
}

func (f *Fetcher) fetchAllBoards(cards []trello.Card) (map[string]*trello.Board, error) {
	boardsChannel := make(chan *trello.Board)
	errorsChannel := make(chan error)

	uniqueBoardIDs := make(map[string]interface{})
	for i := range cards {
		uniqueBoardIDs[cards[i].BoardID] = nil
	}

	for boardID := range uniqueBoardIDs {
		go f.fetchBoard(boardID, boardsChannel, errorsChannel)
	}

	boardsByID := make(map[string]*trello.Board)

	for range uniqueBoardIDs {
		select {
		case board := <-boardsChannel:
			boardsByID[board.ID] = board
		case err := <-errorsChannel:
			return nil, err
		}
	}

	return boardsByID, nil
}

func (f *Fetcher) fetchBoard(boardID string, boardsChannel chan *trello.Board, errorsChannel chan error) {
	board, err := f.Client.GetBoard(boardID)
	if err != nil {
		errorsChannel <- err
		return
	}

	boardsChannel <- board
}

func getProjectBoardID(projectCard *trello.Card) (string, error) {
	boardIDRegex, err := regexp.Compile(regexp.QuoteMeta(boardsURLPath) + `(\w+).*`)
	if err != nil {
		return "", err
	}

	matches := boardIDRegex.FindStringSubmatch(projectCard.Description)
	if len(matches) != 2 {
		return "", fmt.Errorf("could not parse board ID from description on card %s", projectCard.Name)
	}
	return matches[1], nil
}

func getTodoList(lists []trello.List) (*trello.List, error) {
	for _, list := range lists {
		if list.Name == "Todo" {
			return &list, nil
		}
	}
	return nil, fmt.Errorf("missing Todo list on board")
}

func cardsToActions(cards []trello.Card, boardsByID map[string]*trello.Board) []Action {
	actions := make([]Action, 0)
	for i := range cards {
		card := &cards[i]
		actions = append(actions, Action{
			ID:       card.ID,
			Name:     card.Name,
			DueBy:    card.DueBy,
			URL:      card.URL,
			ImageURL: boardsByID[card.BoardID].Preferences.BackgroundImages[0].URL,
		})
	}
	return actions
}
