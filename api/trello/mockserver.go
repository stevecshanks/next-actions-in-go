package trello

import (
	"io/ioutil"
	"net/url"
	"path"

	"github.com/jarcoal/httpmock"
)

// MockServer allows configuring mock responses from a Trello server
type MockServer struct {
	BaseURL *url.URL
	Key     string
	Token   string
}

// CreateMockServer will create and activate a mock server
func CreateMockServer(baseURL, key, token string) *MockServer {
	httpmock.Activate()
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}
	return &MockServer{parsedURL, key, token}
}

// TeardownMockServer will remove the mock server so HTTP responses will behave normally
func TeardownMockServer() {
	defer httpmock.DeactivateAndReset()
}

// AddFileResponse will return the contents of the specified file when the specified path on the mock server is
// requested
func (m *MockServer) AddFileResponse(urlPath, filePath string) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	relativeURL, err := url.Parse(urlPath)
	if err != nil {
		panic(err)
	}

	queryParameters := relativeURL.Query()
	queryParameters.Add("key", m.Key)
	queryParameters.Add("token", m.Token)

	fullURL := m.BaseURL.ResolveReference(&url.URL{
		Path: path.Join(m.BaseURL.Path, relativeURL.Path),
	})

	httpmock.RegisterResponderWithQuery(
		"GET",
		fullURL.String(),
		queryParameters.Encode(),
		httpmock.NewBytesResponder(200, bytes),
	)
}
