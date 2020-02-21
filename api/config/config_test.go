package config

import (
	"fmt"
	"net/url"
	"testing"
)

func TestFromEnvironmentReturnsValidConfig(t *testing.T) {
	SetupEnvironment("http://some.url", "some key", "some token")
	defer TeardownEnvironment()

	config, err := FromEnvironment()
	if err != nil {
		t.Errorf("Error returned from FromEnvironment: %s", err)
	}
	baseURL, _ := url.Parse("http://some.url")
	if config.TrelloBaseURL.String() != baseURL.String() || config.TrelloKey != "some key" || config.TrelloToken != "some token" {
		t.Errorf(fmt.Sprintf("Incorrect config returned from FromEnvironment: %+v", config))
	}
}

func TestFromEnvironmentReturnsErrorIfTrelloBaseURLIsInvalid(t *testing.T) {
	SetupEnvironment(":not a url", "a key", "a token")
	defer TeardownEnvironment()

	_, err := FromEnvironment()
	if err == nil {
		t.Errorf("FromEnvironent did not fail with invalid TRELLO_BASE_URL: %s", err)
	}
}

func TestFromEnvironmentRequiresTrelloBaseURL(t *testing.T) {
	SetupEnvironment("", "a key", "a token")
	defer TeardownEnvironment()

	_, err := FromEnvironment()
	if err == nil {
		t.Errorf("FromEnvironent did not fail with missing TRELLO_BASE_URL: %s", err)
	}
}

func TestFromEnvironmentRequiresTrelloKey(t *testing.T) {
	SetupEnvironment("http://some.url", "", "a token")
	defer TeardownEnvironment()

	_, err := FromEnvironment()
	if err == nil {
		t.Errorf("FromEnvironent did not fail with missing TRELLO_KEY: %s", err)
	}
}

func TestFromEnvironmentRequiresTrelloToken(t *testing.T) {
	SetupEnvironment("http://some.url", "a key", "")
	defer TeardownEnvironment()

	_, err := FromEnvironment()
	if err == nil {
		t.Errorf("FromEnvironent did not fail with missing TRELLO_TOKEN: %s", err)
	}
}
