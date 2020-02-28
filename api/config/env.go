package config

import "os"

// SetupEnvironment sets environment variables specified in the function
func SetupEnvironment(trelloBaseURL string, trelloKey string, trelloToken string, trelloNextActionsListID string) {
	os.Setenv("TRELLO_BASE_URL", trelloBaseURL)
	os.Setenv("TRELLO_KEY", trelloKey)
	os.Setenv("TRELLO_TOKEN", trelloToken)
	os.Setenv("TRELLO_NEXT_ACTIONS_LIST_ID", trelloNextActionsListID)
}

// TeardownEnvironment sets environment variables back to normal
func TeardownEnvironment() {
	os.Setenv("TRELLO_BASE_URL", "")
	os.Setenv("TRELLO_KEY", "")
	os.Setenv("TRELLO_TOKEN", "")
}
