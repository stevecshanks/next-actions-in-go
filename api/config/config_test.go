package config

import (
	"fmt"
	"net/url"
	"testing"
)

func TestFromEnvironmentReturnsValidConfig(t *testing.T) {
	SetupEnvironment("http://some.url", "some key", "some token", "a list id")
	defer TeardownEnvironment()

	config, err := FromEnvironment()
	if err != nil {
		t.Errorf("Error returned from FromEnvironment: %s", err)
	}
	baseURL, _ := url.Parse("http://some.url")
	if config.TrelloBaseURL.String() != baseURL.String() || config.TrelloKey != "some key" || config.TrelloToken != "some token" || config.TrelloNextActionsListID != "a list id" {
		t.Errorf(fmt.Sprintf("Incorrect config returned from FromEnvironment: %+v", config))
	}
}

func TestFromEnvironmentReturnsErrorIfTrelloBaseURLIsInvalid(t *testing.T) {
	SetupEnvironment(":not a url", "a key", "a token", "a list id")
	defer TeardownEnvironment()

	_, err := FromEnvironment()
	if err == nil {
		t.Errorf("FromEnvironment did not fail with invalid TRELLO_BASE_URL: %s", err)
	}
}

func TestFromEnvironmentRequiresTrelloBaseURL(t *testing.T) {
	SetupEnvironment("", "a key", "a token", "a list id")
	defer TeardownEnvironment()

	_, err := FromEnvironment()
	if err == nil {
		t.Errorf("FromEnvironment did not fail with missing TRELLO_BASE_URL: %s", err)
	}
}

func TestFromEnvironmentRequiresTrelloKey(t *testing.T) {
	SetupEnvironment("http://some.url", "", "a token", "a list id")
	defer TeardownEnvironment()

	_, err := FromEnvironment()
	if err == nil {
		t.Errorf("FromEnvironment did not fail with missing TRELLO_KEY: %s", err)
	}
}

func TestFromEnvironmentRequiresTrelloToken(t *testing.T) {
	SetupEnvironment("http://some.url", "a key", "", "a list id")
	defer TeardownEnvironment()

	_, err := FromEnvironment()
	if err == nil {
		t.Errorf("FromEnvironment did not fail with missing TRELLO_TOKEN: %s", err)
	}
}

func TestFromEnvironmentRequiresTrelloNextActionsListID(t *testing.T) {
	SetupEnvironment("http://some.url", "a key", "a token", "")
	defer TeardownEnvironment()

	_, err := FromEnvironment()
	if err == nil {
		t.Errorf("FromEnvironment did not fail with missing TRELLO_NEXT_ACTIONS_LIST_ID: %s", err)
	}
}
