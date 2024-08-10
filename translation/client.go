package translation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

var _ HelloClient = &APIClient{}

type APIClient struct {
	endpoint string
}

// NewHelloClient creates instance of client with a given endpoint
func NewHelloClient(endpoint string) *APIClient {
	return &APIClient{
		endpoint: endpoint,
	}
}

// Translate will call external client for translation.
func (c *APIClient) Translate(word, language string) (string, error) {
	// Prepare query parameters
	params := url.Values{}
	params.Add("language", language)

	// Construct the full URL with query parameters
	fullURL := fmt.Sprintf("%s/%s?%s", c.endpoint, word, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		log.Println(err)
		return "", errors.New("call to api failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		log.Println("word not found")
		return "", nil
	}

	if resp.StatusCode == http.StatusInternalServerError {
		log.Println("error in api")
		return "", errors.New("error in api")
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", errors.New("error reading response")
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return "", errors.New("unable to decode message")
	}

	return m["translation"].(string), nil
}
