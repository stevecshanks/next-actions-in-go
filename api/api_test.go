package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/stevecshanks/next-actions-in-go/api/config"
)

func TestActions(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	bytes, err := ioutil.ReadFile("./trello/testdata/my_cards_response.json")
	if err != nil {
		t.Fatal(err)
	}
	httpmock.RegisterResponderWithQuery(
		"GET",
		"https://api.trello.com/1/members/me/cards",
		"key=some+key&token=some+token",
		httpmock.NewBytesResponder(200, bytes),
	)

	bytes, err = ioutil.ReadFile("./trello/testdata/next_actions_list_response.json")
	if err != nil {
		t.Fatal(err)
	}
	httpmock.RegisterResponderWithQuery(
		"GET",
		"https://api.trello.com/1/lists/nextActionsList123/cards",
		"key=some+key&token=some+token",
		httpmock.NewBytesResponder(200, bytes),
	)

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
