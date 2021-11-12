package config

import (
	"reflect"
	"testing"
)

func TestNewDefaultConfig(t *testing.T) {
	got := newDefaultConfig("path", "root")

	if len(got.Locations) != 1 {
		t.Fatalf("expected exactly 1 location, got %d", len(got.Locations))
	}

	expected := &AppConfig{Locations: []Location{
		{Path: "path", Root: "root"},
	}}

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestLoadConfig(t *testing.T) {
	t.Run("1 configuration file", func(t *testing.T) {
		got := loadConfig("path", "root", []string{"testdata/config/data/config.yml"})

		if got == nil {
			t.Fatalf("want AppConfig, got nil")
		}

		expected := &AppConfig{
			Locations: []Location{
				{"/notes/", "/path/to/notes"},
				{"/tasks/", "/path/to/tasks"},
			}}

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})

	t.Run("2 configuration files", func(t *testing.T) {
		got := loadConfig("path", "root",
			[]string{
				"testdata/config/data/config.yml",
				"testdata/config/home/config.yml",
			})

		if got == nil {
			t.Fatalf("want AppConfig, got nil")
		}

		expected := &AppConfig{
			Locations: []Location{
				{"/notes/", "/path/to/notes"},
				{"/tasks/", "/path/to/tasks"},
				{"/faq/", "/path/to/faq"},
			}}

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})

	t.Run("empty configuration file", func(t *testing.T) {
		got := loadConfig("path", "root",
			[]string{
				"testdata/config/home/err.yml",
			})

		if got == nil {
			t.Fatalf("want AppConfig, got nil")
		}

		expected := newDefaultConfig("path", "root")

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})

	t.Run("malformed configuration file", func(t *testing.T) {
		got := loadConfig("path", "root",
			[]string{
				"testdata/config/Config.YML",
			})

		if got == nil {
			t.Fatalf("want AppConfig, got nil")
		}

		expected := newDefaultConfig("path", "root")

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
}
