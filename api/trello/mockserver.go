package trello

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// CreateMockServer creates a mock Trello server that can be accessed via an HTTP client
func CreateMockServer(t *testing.T, key string, token string) (*httptest.Server, func()) {
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != OwnerCardsPath {
			t.Errorf("Mock server called with incorrect URL path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("key") != key {
			t.Errorf("Mock server called with incorrect key: %s", r.URL.Query().Get("key"))
		}
		if r.URL.Query().Get("token") != token {
			t.Errorf("Mock server called with incorrect token: %s", r.URL.Query().Get("token"))
		}

		json, err := ioutil.ReadFile(filepath.Join(os.Getenv("GOPATH"), "src/next-actions/api/trello/testdata/my_cards_response.json"))
		if err != nil {
			t.Fatal(err)
		}

		w.Write(json)
	})

	server := httptest.NewServer(mockHandler)

	return server, func() {
		server.Close()
	}
}
