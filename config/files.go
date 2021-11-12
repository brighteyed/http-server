package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/OpenPeeDeeP/xdg"
)

const Vendor = ""

// configFiles returns an array of application config files
func configFiles(dirs []string) []string {
	var appConfigFiles []string
	for _, dir := range dirs {
		files, err := findConfigFiles(dir)
		if err == nil {
			appConfigFiles = append(appConfigFiles, files...)
		}
	}

	return appConfigFiles
}

// configDirs returns application config directories
func configDirs() []string {
	xdginfo := xdg.New(Vendor, "http-server")

	return append(
		[]string{xdginfo.ConfigHome()},
		xdginfo.DataDirs()...)
}

// findConfigFiles returns an array of application config files
func findConfigFiles(root string) ([]string, error) {
	var userConfigFiles []string
	fileInfos, err := ioutil.ReadDir(root)
	if err != nil {
		return userConfigFiles, err
	}

	for _, file := range fileInfos {
		if !file.IsDir() {
			if strings.ToLower(filepath.Ext(file.Name())) == ".yml" {
				userConfigFiles = append(userConfigFiles, filepath.Join(root, file.Name()))
			}
		} else {
			contents, err := findConfigFiles(filepath.Join(root, file.Name()))
			if err == nil {
				userConfigFiles = append(userConfigFiles, contents...)
			}
		}
	}

	return userConfigFiles, nil
}
