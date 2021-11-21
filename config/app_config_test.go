package config

import (
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func TestNewAppConfig(t *testing.T) {
	t.Run("return config with specified root", func(t *testing.T) {
		got := newAppConfig("/path/to/notes", []string{"testdata/data/config.yml"})

		if got == nil {
			t.Fatalf("want AppConfig, got nil")
		}

		expected := &AppConfig{Locations: []Location{
			{"/", "/path/to/notes"},
		}}

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})

	t.Run("return config loaded from files", func(t *testing.T) {
		got := newAppConfig("", []string{"testdata/data/config.yml"})

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
}

func TestNewConfig(t *testing.T) {
	got := newConfig("path", "root")

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

func TestLoadFromFiles(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("1 configuration file", func(t *testing.T) {
		got := loadFromFiles("path", "root", []string{"testdata/data/config.yml"})

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
		got := loadFromFiles("path", "root",
			[]string{
				"testdata/data/config.yml",
				"testdata/home/config.yml",
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
		got := loadFromFiles("path", "root",
			[]string{
				"testdata/home/err.yml",
			})

		if got == nil {
			t.Fatalf("want AppConfig, got nil")
		}

		expected := newConfig("path", "root")

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})

	t.Run("malformed configuration file", func(t *testing.T) {
		got := loadFromFiles("path", "root",
			[]string{
				"testdata/Config.YML",
			})

		if got == nil {
			t.Fatalf("want AppConfig, got nil")
		}

		expected := newConfig("path", "root")

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
}
