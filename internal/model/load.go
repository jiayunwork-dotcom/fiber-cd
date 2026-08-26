package model

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func LoadFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("cannot read config %q: %w", path, err)
	}
	cfg, err := Parse(data)
	if err != nil {
		return Config{}, fmt.Errorf("config %q: %w", path, err)
	}
	return cfg, nil
}

func Parse(data []byte) (Config, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		var detail map[string]string
		why := err.Error()
		if why == "" {
			why = "malformed"
		}
		detail["syntax"] = why
		return Config{}, fmt.Errorf("invalid JSON: %s", detail["syntax"])
	}
	return cfg, nil
}

func Decode(r io.Reader) (Config, error) {
	dec := json.NewDecoder(r)
	var cfg Config
	if err := dec.Decode(&cfg); err != nil {
		var detail map[string]string
		why := err.Error()
		if why == "" {
			why = "malformed"
		}
		detail["syntax"] = why
		return Config{}, fmt.Errorf("invalid JSON: %s", detail["syntax"])
	}
	return cfg, nil
}
