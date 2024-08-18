//go:build bdd

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	"github.com/kazakh-in-nz/hello-api/config"
	"github.com/kazakh-in-nz/hello-api/handlers/rest"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
)

const containerName string = "bdd-test"

var opts = godog.Options{
	Output:      colors.Colored(os.Stdout),
	Concurrency: 4,
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestFeatures(t *testing.T) {
	o := opts
	o.TestingT = t

	status := godog.TestSuite{
		Name:                 "godogs",
		Options:              &o,
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
	}.Run()

	if status == 2 {
		t.SkipNow()
	}

	if status != 0 {
		t.Fatalf("zero status code expected, %d received", status)
	}
}

type apiFeature struct {
	client   *resty.Client
	server   *httptest.Server
	word     string
	language string
}

func (api *apiFeature) iTranslateItTo(arg1 string) error {
	api.language = arg1
	return nil
}

func (api *apiFeature) theWord(arg1 string) error {
	api.word = arg1
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	client := resty.New()
	api := &apiFeature{
		client: client,
	}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		cfg := config.Configuration{}
		cfg.LoadFromEnv()

		cfg.DatabaseURL = "localhost"
		cfg.DatabasePort = "6379"

		mux := API(cfg)
		server := httptest.NewServer(mux)

		api.server = server
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		api.server.Close()
		return ctx, nil
	})

	ctx.Step(`^I translate it to "([^"]*)"$`, api.iTranslateItTo)
	ctx.Step(`^the response should be "([^"]*)"$`, api.theResponseShouldBe)
	ctx.Step(`^the word "([^"]*)"$`, api.theWord)
}

func (api *apiFeature) theResponseShouldBe(arg1 string) error {
	url := fmt.Sprintf("%s/%s?language=%s", api.server.URL, api.word, api.language)

	resp, err := api.client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"language": api.language,
		}).
		SetResult(&rest.Resp{}).
		Get(url)

	if err != nil {
		return err
	}

	res := resp.Result().(*rest.Resp)
	if res.Translation != arg1 {
		return fmt.Errorf("translation should be set to %s", arg1)
	}

	return nil
}

var (
	pool     *dockertest.Pool
	database *dockertest.Resource
)

func InitializeTestSuite(sc *godog.TestSuiteContext) {
	var err error

	sc.BeforeSuite(func() {
		dockerSocket := os.Getenv("DOCKER_SOCKET_PATH")
		if dockerSocket == "" {
			osUser := os.Getenv("USER")

			if osUser != "runner" {
				dockerSocket = fmt.Sprintf("unix:///home/%s/.docker/desktop/docker.sock", osUser)
			} else {
				dockerSocket = "unix:///var/run/docker.sock"
			}
		}

		pool, err = dockertest.NewPool(dockerSocket)
		if err != nil {
			panic(fmt.Sprintf("unable to create connection pool %s", err))
		}

		redis, err := pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "redis",
			Tag:        "alpine",
			PortBindings: map[dc.Port][]dc.PortBinding{
				"6379/tcp": {
					{
						HostIP:   "0.0.0.0", // Binding to all interfaces
						HostPort: "6379",    // The port on the host machine
					},
				},
			},
			Name: containerName,
		})
		if err != nil {
			panic(fmt.Sprintf("unable to create container: %s", err))
		}

		if err := redis.Expire(600); err != nil {
			panic("unable to set expiration on container")
		} //Destroy container if it takes too long

		database = redis

		if err := pool.Retry(func() error {
			client := redisClient()
			return client.Ping().Err()
		}); err != nil {
			panic(fmt.Sprintf("could not connect to redis: %s", err))
		}

		// Add key-value pairs to Redis
		client := redisClient()
		err = client.Set("hello:german", "hallo", 0).Err()
		if err != nil {
			panic(fmt.Sprintf("unable to set hello:german in redis: %s", err))
		}

		err = client.Set("hello:bulgarian", "здравейте", 0).Err()
		if err != nil {
			panic(fmt.Sprintf("unable to set hello:bulgarian in redis: %s", err))
		}

		err = client.Set("hello:czech", "ahoj", 0).Err()
		if err != nil {
			panic(fmt.Sprintf("unable to set hello:czech in redis: %s", err))
		}
	})

	sc.AfterSuite(func() {
		database.Close()
		pool.RemoveContainerByName(containerName)
	})
}

func redisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
