package cmd

import (
	"bytes"
	"errors"
	"os"
	"path"
	"testing"
)

func Test_Parse(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dir := path.Join(cwd, "test")
	in := "test/input/version_1"

	t.Run("Unexistable input directory", func(t *testing.T) {
		parser := NewParser()
		p := path.Join(dir, "unknown_folder")
		_, err := parser.Parse(p, nil)

		if err == nil {
			t.Error("Unknown directory should throw error")
		}
	})

	t.Run("Getwd() error bubbling", func(t *testing.T) {
		parser := &Parser{
			Wd: func() (string, error) {
				return "", errors.New("WD ERROR")
			},
		}
		_, err := parser.Parse("some/folder/", nil)

		if err == nil || err.Error() != "WD ERROR" {
			t.Error("Getwd() error should return error")
		}
	})

	t.Run("Unexistable output directory", func(t *testing.T) {
		parser := NewParser()
		unexisting := path.Join(dir, "input", "unexisting")
		_, err := parser.Parse(unexisting, nil)

		if err == nil {
			t.Error("Unknown directory should return error")
		}
	})

	t.Run("Empty input directory", func(t *testing.T) {
		parser := NewParser()
		p := path.Join(dir, "input", "version_2")
		_, err := parser.Parse(p, nil)

		if err == nil || err.Error() != errEmptyDirectory {
			t.Error("Empty directory should return error")
		}
	})

	t.Run("Empty result from parsing", func(t *testing.T) {
		parser := NewParser()
		p := path.Join(dir, "input", "version_3")
		_, err := parser.Parse(p, nil)

		if err == nil || err.Error() != errEmptyResult {
			t.Error("Empty result from parsing should return error")
		}
	})

	t.Run("Error on wrong template", func(t *testing.T) {
		parser := NewParser()
		p := path.Join(dir, "input", "version_4")
		_, err := parser.Parse(p, nil)

		if err == nil {
			t.Error("Wrong template should return error")
		}
	})

	t.Run("Parse template successfully", func(t *testing.T) {
		parser := NewParser()
		buf, err := parser.Parse(in, map[string]string{
			"Title": "Testing",
		})

		if err != nil {
			t.Error(err)
		}

		if !bytes.Contains(buf, []byte("# Testing")) {
			t.Error("Template don't contain env variables")
		}
	})
}
