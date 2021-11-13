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

// NewAppConfig returns server configuration that is defined by root or
// (if root is empty) loaded from configuration files
func NewAppConfig(root string) *AppConfig {
	return newAppConfig(root, configFiles(configDirs()))
}

// newConfig loads config from files if root is empty. Otherwise it  returns
// config with one location defined by root
func newAppConfig(root string, configFiles []string) *AppConfig {
	var appCfg *AppConfig
	if root == "" {
		appCfg = loadFromFiles("/", ".", configFiles)
	} else {
		appCfg = newConfig("/", root)
	}
	return appCfg
}

// newConfig returns application configuration with one location
func newConfig(path string, root string) *AppConfig {
	return &AppConfig{Locations: []Location{
		{Path: path, Root: root},
	}}
}

// loadFromFiles merges all locations from configuration files. If no locations were loaded
// it returns location with specified defaultRoot and defaultPath
func loadFromFiles(defaultPath string, defaultRoot string, configFiles []string) *AppConfig {
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
		return newConfig(defaultPath, defaultRoot)
	}

	return &cfg
}
