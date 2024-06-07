package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestFindConfigFiles(t *testing.T) {
	got, err := findConfigFiles("testdata")

	if err != nil {
		t.Fatalf("got an unexpected error")
	}

	if len(got) != 5 {
		t.Errorf("expected 5 files, got %d", len(got))
	}

	for _, file := range got {
		if strings.ToLower(filepath.Ext(file)) != ".yml" {
			t.Errorf("%q: expected %q file extension", file, ".yml")
		}
	}
}

func TestConfigFiles(t *testing.T) {
	got := configFiles([]string{"testdata/data", "testdata/home"})

	if len(got) != 4 {
		t.Errorf("expected 4 files, got %d", len(got))
	}
}

func TestConfigDirs(t *testing.T) {
	got := configDirs()

	if len(got) < 2 {
		t.Errorf("expected at least 2 config directories but found only %d", len(got))
	}
}
