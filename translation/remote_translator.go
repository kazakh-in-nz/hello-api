package translation

import (
	"fmt"
	"strings"
)

type RemoteService struct {
	client HelloClient
	cache  map[string]string
}

type HelloClient interface {
	Translate(world, language string) (string, error)
}

func NewRemoteService(client HelloClient) *RemoteService {
	return &RemoteService{
		client: client,
		cache:  make(map[string]string),
	}
}

func (s *RemoteService) Translate(word string, language string) string {
	word = strings.ToLower(word)
	language = strings.ToLower(language)

	key := fmt.Sprintf("%s:%s", word, language)

	tr, ok := s.cache[key]
	if ok {
		return tr
	}

	resp, err := s.client.Translate(word, language)
	if err != nil {
		return ""
	}
	s.cache[key] = resp
	return resp
}
