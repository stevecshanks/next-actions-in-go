package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
	mockServer := trello.CreateMockServer("some key", "some token")
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
	mockServer.AddFileResponse(
		trello.BoardPath("myBoardId"),
		trelloResponse("board_response.json"),
	)
	mockServer.AddFileResponse(
		trello.BoardPath("boardWithNoImagesId"),
		trelloResponse("board_with_no_images_response.json"),
	)

	config.SetupEnvironment("some key", "some token", "nextActionsList123", "projectsList456")
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

	assertResponseMatchesContractFile(t, rr.Body.Bytes(), "api_success_response.json")
}

func TestActionsErrors(t *testing.T) {
	trello.CreateMockServer("some key", "some token")
	defer trello.TeardownMockServer()

	config.SetupEnvironment("some key", "some token", "nextActionsList123", "projectsList456")
	defer config.TeardownEnvironment()

	req, err := http.NewRequest("GET", "/actions", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(actions)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("/actions returned status: %v", status)
	}

	assertResponseMatchesContractFile(t, rr.Body.Bytes(), "api_error_response.json")
}

func assertResponseMatchesContractFile(t *testing.T, response []byte, fileName string) {
	expectedBytes, err := ioutil.ReadFile(path.Join("../../../contracts", fileName))
	if err != nil {
		panic(err)
	}

	var actualBytes bytes.Buffer
	err = json.Indent(&actualBytes, response, "", "  ")
	if err != nil {
		t.Fatalf("Could not parse response as JSON")
	}

	if !bytes.Equal(actualBytes.Bytes(), expectedBytes) {
		t.Errorf(
			"Response did not match file %s:\nexpected: %v\nactual:   %v",
			fileName,
			string(expectedBytes),
			actualBytes.String(),
		)
	}
}
