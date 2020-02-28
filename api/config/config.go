package config

import (
	"fmt"
	"net/url"
	"os"
)

// Config represents a configuration for the app
type Config struct {
	TrelloBaseURL           *url.URL
	TrelloKey               string
	TrelloToken             string
	TrelloNextActionsListID string
}

// FromEnvironment creates a Config from environment variables
func FromEnvironment() (*Config, error) {
	trelloBaseURLString, err := requiredEnvironmentVariable("TRELLO_BASE_URL")
	if err != nil {
		return nil, err
	}
	trelloBaseURL, err := url.Parse(trelloBaseURLString)
	if err != nil {
		return nil, err
	}

	trelloKey, err := requiredEnvironmentVariable("TRELLO_KEY")
	if err != nil {
		return nil, err
	}

	trelloToken, err := requiredEnvironmentVariable("TRELLO_TOKEN")
	if err != nil {
		return nil, err
	}

	trelloNextActionsListID, err := requiredEnvironmentVariable("TRELLO_NEXT_ACTIONS_LIST_ID")
	if err != nil {
		return nil, err
	}

	return &Config{
		TrelloBaseURL:           trelloBaseURL,
		TrelloKey:               trelloKey,
		TrelloToken:             trelloToken,
		TrelloNextActionsListID: trelloNextActionsListID,
	}, nil
}

func requiredEnvironmentVariable(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("%s is a required environment variable", name)
	}
	return value, nil
}
