package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stevecshanks/next-actions-in-go/api/config"
	"github.com/stevecshanks/next-actions-in-go/api/trello"
)

func TestActions(t *testing.T) {
	mockServer := trello.CreateMockServer("https://api.trello.com/1", "some key", "some token")
	defer trello.TeardownMockServer()

	mockServer.AddFileResponse("members/me/cards", "./trello/testdata/my_cards_response.json")
	mockServer.AddFileResponse("lists/nextActionsList123/cards", "./trello/testdata/next_actions_list_response.json")

	config.SetupEnvironment("https://api.trello.com/1", "some key", "some token", "nextActionsList123")
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
		`{"type":"actions","id":"111111111111111111111111","name":"My First Card"},` +
		`{"type":"actions","id":"222222222222222222222222","name":"My Second Card"},` +
		`{"type":"actions","id":"333333333333333333333333","name":"My Third Card"}` +
		`]}` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("/actions returned incorrect body:\nexpected: %v\nactual:   %v", expected, rr.Body.String())
	}
}
