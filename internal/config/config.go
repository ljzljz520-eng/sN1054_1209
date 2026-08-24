package config

import "strings"

type Config struct {
	DatabasePath string
	Address      string
	ReadOnly     bool
}

func Default() Config { return Config{DatabasePath: "familyitinerary.db", Address: ":8080"} }

func FromValues(databasePath, address string, readOnly bool) Config {
	config := Default()
	if strings.TrimSpace(databasePath) != "" {
		config.DatabasePath = databasePath
	}
	if strings.TrimSpace(address) != "" {
		config.Address = address
	}
	config.ReadOnly = readOnly
	return config
}

func (c Config) Valid() bool { return c.DatabasePath != "" && c.Address != "" }

func (c Config) Endpoint(path string) string {
	if strings.HasPrefix(path, "/") {
		return c.Address + path
	}
	return c.Address + "/" + path
}
