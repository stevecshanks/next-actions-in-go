package trello // nolint:golint // package comment is in another file

// List represents a Trello list returned via the API
type List struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
