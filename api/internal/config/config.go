package config

import (
	"fmt"
	"os"
)

// Config represents a configuration for the app
type Config struct {
	TrelloKey               string
	TrelloToken             string
	TrelloNextActionsListID string
	TrelloProjectsListID    string
}

// FromEnvironment creates a Config from environment variables
func FromEnvironment() (*Config, error) {
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

	trelloProjectsListID, err := requiredEnvironmentVariable("TRELLO_PROJECTS_LIST_ID")
	if err != nil {
		return nil, err
	}

	return &Config{
		TrelloKey:               trelloKey,
		TrelloToken:             trelloToken,
		TrelloNextActionsListID: trelloNextActionsListID,
		TrelloProjectsListID:    trelloProjectsListID,
	}, nil
}

func requiredEnvironmentVariable(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("%s is a required environment variable", name)
	}
	return value, nil
}
