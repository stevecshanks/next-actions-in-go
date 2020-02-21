package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"next-actions/api/config"
	"next-actions/api/trello"
	"testing"
)

func TestActions(t *testing.T) {
	server, teardown := trello.CreateMockServer(t, "some key", "some token")
	defer teardown()

	config.SetupEnvironment(fmt.Sprintf("http://%s", server.Listener.Addr().String()), "some key", "some token")
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
		`{"type":"actions","id":"222222222222222222222222","name":"My Second Card"}` +
		`]}` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("/actions returned incorrect body:\nexpected: %v\nactual:   %v", expected, rr.Body.String())
	}
}
