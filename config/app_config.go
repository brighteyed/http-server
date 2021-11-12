package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Location defines a directory to serve
type Location struct {
	Path string `yaml:"path"`
	Root string `yaml:"root"`
}

// AppConfig defines server configuration
type AppConfig struct {
	Locations []Location `yaml:"locations"`
}

// LoadConfig loads application config from yaml files or uses default args
// as fallback
func LoadConfig(defaultPath string, defaultRoot string) *AppConfig {
	return loadConfig(defaultPath, defaultRoot, configFiles(configDirs()))
}

func loadConfig(defaultPath string, defaultRoot string, configFiles []string) *AppConfig {
	cfg := AppConfig{}

	for _, config := range configFiles {
		content, err := ioutil.ReadFile(config)
		if err != nil {
			log.Printf("Can't read config file %q\n", config)
			continue
		}

		curr := AppConfig{}
		if err := yaml.Unmarshal(content, &curr); err != nil {
			log.Printf("Can't parse config file %q, %v", config, err)
			continue
		}

		cfg.Locations = append(cfg.Locations, curr.Locations...)
	}

	if len(cfg.Locations) == 0 {
		return newDefaultConfig(defaultPath, defaultRoot)
	}

	return &cfg
}

func newDefaultConfig(path string, root string) *AppConfig {
	return &AppConfig{Locations: []Location{
		{Path: path, Root: root},
	}}
}
