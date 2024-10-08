//go:build unit

package translation_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kazakh-in-nz/hello-api/translation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestHelloClientSuite(t *testing.T) {
	suite.Run(t, new(HelloClientSuite))
}

type HelloClientSuite struct {
	suite.Suite
	mockServerService *MockService
	server            *httptest.Server
	underTest         translation.HelloClient
}

func (suite *HelloClientSuite) SetupSuite() {
	suite.mockServerService = new(MockService)

	handler := func(w http.ResponseWriter, r *http.Request) {
		language := r.URL.Query().Get("language")
		word := strings.ReplaceAll(r.URL.Path, "/", "")

		if word == "" || language == "" {
			http.Error(w, "invalid input", 400)
			return
		}

		resp, err := suite.mockServerService.Translate(word, language)
		if err != nil {
			http.Error(w, "error", 500)
			return
		}

		if resp == "" {
			http.Error(w, "missing", 404)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, resp)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	suite.server = httptest.NewServer(mux)

	suite.underTest = translation.NewHelloClient(suite.server.URL)
}

func (suite *HelloClientSuite) SetupTest() {
	suite.mockServerService = new(MockService)
}

func (suite *HelloClientSuite) TearDownSuite() {
	suite.server.Close()
}

type MockService struct {
	mock.Mock
}

func (m *MockService) Translate(word, language string) (string, error) {
	args := m.Called(word, language)
	return args.String(0), args.Error(1)
}

func (suite *HelloClientSuite) TestCall() {
	// Arrange
	suite.mockServerService.On("Translate", "foo", "bar").Return(`{
    "translation":"baz"}`, nil)

	// Act
	resp, err := suite.underTest.Translate("foo", "bar")

	// Assert
	suite.NoError(err)
	suite.Equal(resp, "baz")
}

func (suite *HelloClientSuite) TestCall_APIError() {
	// Arrange
	suite.mockServerService.On("Translate", "foo", "bar").Return("", errors.New("this is a test"))

	// Act
	resp, err := suite.underTest.Translate("foo", "bar")

	// Assert
	suite.EqualError(err, "error in api")
	suite.Equal(resp, "")
}

func (suite *HelloClientSuite) TestCall_InvalidJSON() {
	// Arrange
	suite.mockServerService.On("Translate", "foo", "bar").Return(`invalid 
    json`, nil)

	// Act
	resp, err := suite.underTest.Translate("foo", "bar")

	// Assert
	suite.EqualError(err, "unable to decode message")
	suite.Equal(resp, "")
}
