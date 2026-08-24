package config

import "testing"

func TestConfigValues(t *testing.T) {
	config := FromValues("trip.db", ":9090", true)
	if !config.Valid() || config.DatabasePath != "trip.db" || config.Address != ":9090" || !config.ReadOnly {
		t.Fatal("config values incorrect")
	}
	if config.Endpoint("health") != ":9090/health" || config.Endpoint("/health") != ":9090/health" {
		t.Fatal("endpoint formatting incorrect")
	}
}
