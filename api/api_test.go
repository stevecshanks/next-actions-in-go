package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestActions(t *testing.T) {
	trelloKey := "my_trello_key"
	trelloToken := "my_trello_token"

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/members/me/cards" {
			t.Errorf("Mock server called with incorrect URL path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("key") != trelloKey {
			t.Errorf("Mock server called with incorrect key: %s", r.URL.Query().Get("key"))
		}
		if r.URL.Query().Get("token") != trelloToken {
			t.Errorf("Mock server called with incorrect token: %s", r.URL.Query().Get("token"))
		}

		json, err := ioutil.ReadFile("fixtures/my_cards_response.json")
		if err != nil {
			t.Fatal(err)
		}

		w.Write(json)
	})

	server := httptest.NewServer(mockHandler)
	defer server.Close()

	os.Setenv("TRELLO_BASE_URL", fmt.Sprintf("http://%s", server.Listener.Addr().String()))
	os.Setenv("TRELLO_KEY", trelloKey)
	os.Setenv("TRELLO_TOKEN", trelloToken)

	defer os.Setenv("TRELLO_BASE_URL", "")
	defer os.Setenv("TRELLO_KEY", "")
	defer os.Setenv("TRELLO_TOKEN", "")

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
		`{"type":"actions","id":"222222222222222222222222","name":"My Second Card"}` +
		`]}` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("/actions returned incorrect body:\nexpected: %v\nactual:   %v", expected, rr.Body.String())
	}
}
