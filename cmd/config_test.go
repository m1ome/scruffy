package cmd

import (
	"errors"
	"testing"
)

func TestParseConfig(t *testing.T) {
	t.Run("Parse non existing file", func(t *testing.T) {
		c := NewConfig()
		err := c.Parse("non_existant_file")

		if err == nil {
			t.Error("Unexistant file should return error")
		}
	})

	t.Run("Getwd() error should raise error", func(t *testing.T) {
		c := &Config{
			Wd: func() (string, error) {
				return "", errors.New("Some Bad Error")
			},
		}

		err := c.Parse("test/config.yml")
		if err == nil || err.Error() != "Some Bad Error" {
			t.Error("Getwd() error should raise error")
		}
	})

	t.Run("Parse corrupted file", func(t *testing.T) {
		c := NewConfig()
		err := c.Parse("test/wrong_config.yml")

		if err == nil {
			t.Error("Corrupted file should return error")
		}
	})

	t.Run("Parse configuration file", func(t *testing.T) {
		c := NewConfig()
		err := c.Parse("test/config.yml")

		if err != nil {
			t.Errorf("Config parsing returned: %s", err)
		}

		if c.YML.Environments == nil || len(c.YML.Environments) != 2 {
			t.Error("Parsed wrong number of environments from config")
		}

		env, err := c.Env("public")
		if err != nil {
			t.Error("Cannot get public environment")
		}

		if env.Release != "scruffypublic" {
			t.Error("Public Name parsed wrong")
		}

		_, err = c.Env("wrong_env")
		if err == nil {
			t.Error("Wrong env should return error")
		}
	})
}
