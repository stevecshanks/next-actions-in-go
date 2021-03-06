package config // nolint:golint // package comment is in another file

import "os"

// SetupEnvironment sets environment variables specified in the function
func SetupEnvironment(trelloKey, trelloToken, trelloNextActionsListID, trelloProjectsListID string) {
	os.Setenv("TRELLO_KEY", trelloKey)
	os.Setenv("TRELLO_TOKEN", trelloToken)
	os.Setenv("TRELLO_NEXT_ACTIONS_LIST_ID", trelloNextActionsListID)
	os.Setenv("TRELLO_PROJECTS_LIST_ID", trelloProjectsListID)
}

// TeardownEnvironment sets environment variables back to normal
func TeardownEnvironment() {
	os.Setenv("TRELLO_KEY", "")
	os.Setenv("TRELLO_TOKEN", "")
	os.Setenv("TRELLO_NEXT_ACTIONS_LIST_ID", "")
	os.Setenv("TRELLO_PROJECTS_LIST_ID", "")
}
