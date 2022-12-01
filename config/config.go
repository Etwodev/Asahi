package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const CONFIG = "./config.json"

var c *Config

type Config struct {
	port         string
	address      string
	assets		 string
	experimental bool
}

func (c Config) Port() string {
	return c.port
}

func (c Config) Address() string {
	return c.address
}

func (c Config) Experimental() bool {
	return c.experimental
}

func (c Config) Assets() string {
	return c.assets
}


func Load() error {
	_, err := os.Stat(CONFIG)
	if os.IsNotExist(err) {
		if err := Create(); err != nil {
			return fmt.Errorf("Load: failed creating load: %w", err)
		}
	}

	file, err := ioutil.ReadFile(CONFIG)
	if err != nil {
		return fmt.Errorf("Load: failed reading json: %w", err)
	}

	err = json.Unmarshal(file, &c)
	if err != nil {
		return fmt.Errorf("Load: failed marshalling json: %w", err)
	}
	return nil
}

func Create() error {
	file, err := json.MarshalIndent(&Config{port: "8080", address: "localhost", experimental: false, assets: "./assets"}, "", " ")
	if err != nil {
		return fmt.Errorf("Create: failed marshalling config: %w", err)
	}
	err = ioutil.WriteFile(CONFIG, file, 0644)
	if err != nil {
		return fmt.Errorf("Create: failed writing config: %w", err)
	}
	return nil
}

func New() ( *Config, error ) {
	if c == nil {
		err := Load()
		if err != nil {
			return nil, fmt.Errorf("New: failed loading json: %w", err)
		}
	}
	return c, nil
}