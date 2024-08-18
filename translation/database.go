package translation

import (
	"context"
	"fmt"

	"github.com/kazakh-in-nz/hello-api/config"
	"github.com/kazakh-in-nz/hello-api/handlers/rest"
	"github.com/redis/go-redis/v9"
)

var _ rest.Translator = &Database{}

type Database struct {
	conn *redis.Client
}

func NewDatabaseService(cfg config.Configuration) *Database {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.DatabaseURL, cfg.DatabasePort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &Database{
		conn: rdb,
	}
}

func (s *Database) Close() error {
	return s.conn.Close()
}

func (s *Database) Translate(word string, language string) string {
	out := s.conn.Get(context.Background(), fmt.Sprintf("%s:%s",
		word, language))
	return out.Val()
}
