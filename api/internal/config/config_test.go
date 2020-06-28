package config

import (
	"fmt"
	"testing"
)

func TestFromEnvironmentReturnsValidConfig(t *testing.T) {
	SetupEnvironment("some key", "some token", "next actions list id", "projects list id")
	defer TeardownEnvironment()

	config, err := FromEnvironment()
	if err != nil {
		t.Fatalf("Error returned from FromEnvironment: %s", err)
	}

	isValidConfig := config.TrelloKey == "some key" &&
		config.TrelloToken == "some token" &&
		config.TrelloNextActionsListID == "next actions list id" &&
		config.TrelloProjectsListID == "projects list id"

	if !isValidConfig {
		t.Errorf(fmt.Sprintf("Incorrect config returned from FromEnvironment: %+v", config))
	}
}

func TestFromEnvironmentRequiresTrelloKey(t *testing.T) {
	SetupEnvironment("", "a token", "a list id", "projects list id")
	defer TeardownEnvironment()

	_, err := FromEnvironment()
	if err == nil {
		t.Errorf("FromEnvironment did not fail with missing TRELLO_KEY: %s", err)
	}
}

func TestFromEnvironmentRequiresTrelloToken(t *testing.T) {
	SetupEnvironment("a key", "", "next actions list id", "projects list id")
	defer TeardownEnvironment()

	_, err := FromEnvironment()
	if err == nil {
		t.Errorf("FromEnvironment did not fail with missing TRELLO_TOKEN: %s", err)
	}
}

func TestFromEnvironmentRequiresTrelloNextActionsListID(t *testing.T) {
	SetupEnvironment("a key", "a token", "", "projects list id")
	defer TeardownEnvironment()

	_, err := FromEnvironment()
	if err == nil {
		t.Errorf("FromEnvironment did not fail with missing TRELLO_NEXT_ACTIONS_LIST_ID: %s", err)
	}
}

func TestFromEnvironmentRequiresTrelloProjectsListID(t *testing.T) {
	SetupEnvironment("a key", "a token", "next actions list id", "")
	defer TeardownEnvironment()

	_, err := FromEnvironment()
	if err == nil {
		t.Errorf("FromEnvironment did not fail with missing TRELLO_PROJECTS_LIST_ID: %s", err)
	}
}
