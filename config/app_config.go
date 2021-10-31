package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/OpenPeeDeeP/xdg"
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

func newDefaultConfig(path string, root string) *AppConfig {
	return &AppConfig{Locations: []Location{
		{Path: path, Root: root},
	}}
}

// LoadConfig loads application config from yaml or uses default args
// as fallback
func LoadConfig(defaultPath string, defaultRoot string) *AppConfig {
	configFile, err := configFile()
	if err != nil {
		log.Printf("Looking for config file, %v", err)
		return nil
	}

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("Can't read config file. Using default root to serve")
		return newDefaultConfig(defaultPath, defaultRoot)
	}

	cfg := AppConfig{}
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		log.Printf("Can't parse config file. Using default root to serve, %v", err)
		return newDefaultConfig(defaultPath, defaultRoot)
	}

	return &cfg
}

// configFile returns a path to the application config file
func configFile() (string, error) {
	folder, err := findOrCreateConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(folder, "config.yml"), nil
}

// findOrCreateConfigDir ensures application config directory exists
func findOrCreateConfigDir() (string, error) {
	folder := xdg.New("", "http-server").ConfigHome()
	return folder, os.MkdirAll(folder, 0755)
}
