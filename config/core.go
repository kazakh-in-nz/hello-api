package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Configuration struct {
	Port            string `json:"port"`
	DefaultLanguage string `json:"default_language"`
	LegacyEndpoint  string `json:"legacy_endpoint"`
	DatabaseType    string `json:"database_type"`
	DatabaseURL     string `json:"database_url"`
	DatabasePort    string `json:"database_port"`
}

var defaultConfig = Configuration{
	Port:            "8080",
	DefaultLanguage: "english",
}

func (c *Configuration) LoadFromEnv() {
	if lang := os.Getenv("DEFAULT_LANGUAGE"); lang != "" {
		c.DefaultLanguage = lang
	}

	if port := os.Getenv("PORT"); port != "" {
		c.Port = port
	}
}

func (c *Configuration) ParsePort() {
	if c.Port[0] != ':' {
		c.Port = ":" + c.Port
	}

	if _, err := strconv.Atoi(c.Port[1:]); err != nil {
		fmt.Printf("invalid port %s", c.Port)
		c.Port = defaultConfig.Port
	}
}

func (c *Configuration) LoadFromJSON(path string) error {
	log.Printf("loading configuration from file %s", path)

	b, err := os.ReadFile(path)
	if err != nil {
		log.Printf("unable to load file: %s\n", err.Error())
		return errors.New("unable to load configuration")
	}

	if err := json.Unmarshal(b, c); err != nil {
		log.Printf("unable to parse file: %s\n", err.Error())
		return errors.New("unable to load configuraiton")
	}

	if c.Port == "" {
		log.Printf("empty port, reverting to default")
		c.Port = defaultConfig.Port
	}

	if c.DefaultLanguage == "" {
		log.Printf("empty language, reverting to default")
		c.DefaultLanguage = defaultConfig.DefaultLanguage
	}

	return nil
}

// LoadConfiguration will provide cycle through flags, files, and finally env variables to load configuration.
func LoadConfiguration() Configuration {
	cfgfileFlag := flag.String("config_file", "", "load configurations from a file")
	portFlag := flag.String("port", "", "set port")

	flag.Parse()
	cfg := defaultConfig

	if cfgfileFlag != nil && *cfgfileFlag != "" {
		if err := cfg.LoadFromJSON(*cfgfileFlag); err != nil {
			log.Printf("unable to load configuration from json: %s, using default values", *cfgfileFlag)
		}
	}

	cfg.LoadFromEnv()

	if portFlag != nil && *portFlag != "" {
		cfg.Port = *portFlag
	}

	cfg.ParsePort()
	return cfg
}
