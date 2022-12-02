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
	Port         string `json:"port"`
	Address      string `json:"address"`
	Assets		 string	 `json:"assets"`
	Experimental bool	 `json:"experimental"`
}

func Port() string {
	return c.Port
}

func Address() string {
	return c.Address
}

func Experimental() bool {
	return c.Experimental
}

func Assets() string {
	return c.Assets
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
	file, err := json.MarshalIndent(&Config{Port: "8080", Address: "localhost", Experimental: false, Assets: "./assets"}, "", " ")
	if err != nil {
		return fmt.Errorf("Create: failed marshalling config: %w", err)
	}
	err = ioutil.WriteFile(CONFIG, file, 0644)
	if err != nil {
		return fmt.Errorf("Create: failed writing config: %w", err)
	}
	return nil
}

func New() ( error ) {
	if c == nil {
		err := Load()
		if err != nil {
			return fmt.Errorf("New: failed loading json: %w", err)
		}
	}
	return nil
}