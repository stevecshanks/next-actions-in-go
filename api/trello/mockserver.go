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
func CreateMockServer(baseURL string, key string, token string) *MockServer {
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
func (m *MockServer) AddFileResponse(urlPath string, filePath string) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	queryParameters := m.BaseURL.Query()
	queryParameters.Add("key", m.Key)
	queryParameters.Add("token", m.Token)

	httpmock.RegisterResponderWithQuery(
		"GET",
		path.Join(m.BaseURL.Path, urlPath),
		queryParameters.Encode(),
		httpmock.NewBytesResponder(200, bytes),
	)
}
