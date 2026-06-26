package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ThemeConfig struct {
	Primary    string `yaml:"primary"`
	Secondary  string `yaml:"secondary"`
	Background string `yaml:"background"`
}

type BackendConfig struct {
	SwarmAPIURL   string `yaml:"swarmApiUrl"`
	TicketAPIURL  string `yaml:"ticketApiUrl"`
}

type TenantConfig struct {
	AppName string        `yaml:"appName"`
	Theme   ThemeConfig   `yaml:"theme"`
	Backend BackendConfig `yaml:"backend"`
}

func LoadTenantConfig(path string) (TenantConfig, error) {
	var cfg TenantConfig

	f, err := os.Open(path)
	if err != nil {
		// fall back to env-only config
		return loadFromEnv(), nil
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return loadFromEnv(), nil
	}

	// env overrides
	if v := os.Getenv("COCKPIT_APP_NAME"); v != "" {
		cfg.AppName = v
	}
	if v := os.Getenv("COCKPIT_SWARM_API_URL"); v != "" {
		cfg.Backend.SwarmAPIURL = v
	}
	if v := os.Getenv("COCKPIT_TICKET_API_URL"); v != "" {
		cfg.Backend.TicketAPIURL = v
	}

	return cfg, nil
}

func loadFromEnv() TenantConfig {
	return TenantConfig{
		AppName: os.Getenv("COCKPIT_APP_NAME"),
		Theme: ThemeConfig{
			Primary:    "#ff1493",
			Secondary:  "#00ff7f",
			Background: "#000000",
		},
		Backend: BackendConfig{
			SwarmAPIURL:  os.Getenv("COCKPIT_SWARM_API_URL"),
			TicketAPIURL: os.Getenv("COCKPIT_TICKET_API_URL"),
		},
	}
}
