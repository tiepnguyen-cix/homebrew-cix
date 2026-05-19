package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds user-defined overrides loaded from .cix.yml.
type Config struct {
	Variables map[string]string `yaml:"variables"`
	Secrets   map[string]string `yaml:"secrets"`
	Docker    DockerConfig      `yaml:"docker"`
}

// DockerConfig holds Docker-specific settings.
type DockerConfig struct {
	PullPolicy string `yaml:"pull_policy"` // always, if-not-present, never
	Network    string `yaml:"network"`
}

// Load reads .cix.yml from the given path.
// If the file does not exist, an empty Config is returned without error.
func Load(path string) (*Config, error) {
	cfg := &Config{
		Variables: make(map[string]string),
		Secrets:   make(map[string]string),
		Docker: DockerConfig{
			PullPolicy: "if-not-present",
			Network:    "bridge",
		},
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// Resolve secrets from shell environment.
	for k, v := range cfg.Secrets {
		if v == "" {
			continue
		}
		// Strip leading $ if present, then look up env var.
		envKey := v
		if len(envKey) > 0 && envKey[0] == '$' {
			envKey = envKey[1:]
		}
		if envVal := os.Getenv(envKey); envVal != "" {
			cfg.Variables[k] = envVal
		}
	}

	return cfg, nil
}