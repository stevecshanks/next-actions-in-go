package main

import (
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/stevecshanks/next-actions-in-go/api/internal/config"
	"github.com/stevecshanks/next-actions-in-go/api/internal/trello"
)

func trelloResponse(fileName string) string {
	return path.Join("../../internal/trello/testdata", fileName)
}

func TestActions(t *testing.T) {
	mockServer := trello.CreateMockServer("https://api.trello.com/1", "some key", "some token")
	defer trello.TeardownMockServer()

	mockServer.AddFileResponse(trello.OwnedCardsPath(), trelloResponse("my_cards_response.json"))
	mockServer.AddFileResponse(
		trello.CardsOnListPath("nextActionsList123"),
		trelloResponse("next_actions_list_response.json"),
	)
	mockServer.AddFileResponse(
		trello.CardsOnListPath("projectsList456"),
		trelloResponse("projects_list_response.json"),
	)
	mockServer.AddFileResponse(
		trello.ListsOnBoardPath("projectBoard789"),
		trelloResponse("board_lists_response.json"),
	)
	mockServer.AddFileResponse(
		trello.CardsOnListPath("todoListId"),
		trelloResponse("project_todo_list_cards_response.json"),
	)

	config.SetupEnvironment("https://api.trello.com/1", "some key", "some token", "nextActionsList123", "projectsList456")
	defer config.TeardownEnvironment()

	req, err := http.NewRequest("GET", "/actions", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(actions)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("/actions returned status: %v", status)
	}

	expected := `{"data":[` +
		`{"type":"actions","id":"myFirstCardId","name":"My First Card","dueBy":"2020-02-12T16:24:00Z",` +
		`"url":"https://trello.com/c/abcd1234/10-my-first-card"},` +
		`{"type":"actions","id":"mySecondCardId","name":"My Second Card","dueBy":null,` +
		`"url":"https://trello.com/c/bcde2345/11-my-second-card"},` +
		`{"type":"actions","id":"todoCardId","name":"Todo Card","dueBy":null,` +
		`"url":"https://trello.com/c/cdef3456/33-my-third-card"},` +
		`{"type":"actions","id":"firstProjectCardId","name":"Project Card","dueBy":null,` +
		`"url":"https://trello.com/c/fghi5678/55-my-project-card"}` +
		`]}` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("/actions returned incorrect body:\nexpected: %v\nactual:   %v", expected, rr.Body.String())
	}
}
