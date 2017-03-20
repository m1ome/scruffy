package cmd

import (
	"bytes"
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
	in := path.Join(dir, "input", "version_1")

	t.Run("Unexistable input directory", func(t *testing.T) {
		p := path.Join(dir, "unknown_folder")
		_, err := Parse(p, nil)

		if err == nil {
			t.Error("Unknown directory should throw error")
		}
	})

	t.Run("Unexistable output directory", func(t *testing.T) {
		_, err := Parse(in, nil)

		if err == nil {
			t.Error("Unknown directory should return error")
		}
	})

	t.Run("Empty input directory", func(t *testing.T) {
		p := path.Join(dir, "input", "version_2")
		_, err := Parse(p, nil)

		if err == nil || err.Error() != errEmptyDirectory {
			t.Error("Empty directory should return error")
		}
	})

	t.Run("Empty result from parsing", func(t *testing.T) {
		p := path.Join(dir, "input", "version_3")
		_, err := Parse(p, nil)

		if err == nil || err.Error() != errEmptyResult {
			t.Error("Empty result from parsing should return error")
		}
	})

	tpl := []byte(`
# Testing

## Accounts [/accounts]

### List [GET]

+ Response (text/plain)

        Hello, world!

+ Response (application/json)

        {"message": "Hello, world!"}

### Delete [/accounts/{id} GET]
+ Parameters
    + id (number, required) - Account ID
`)

	t.Run("Parse template successfully", func(t *testing.T) {
		buf, err := Parse(in, map[string]string{
			"Title": "Testing",
		})

		if err != nil {
			t.Error(err)
		}

		if bytes.Compare(buf, tpl) != 0 {
			t.Error("Templates differs")
		}
	})
}
