package cmd

import (
	"fmt"
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

// Config YAML configuration object handler
type Config struct {
	YML *configYML
	Wd  WorkingDirGetter
}

type configYML struct {
	Token        string         `yaml:"token"`
	Source       string         `yaml:"source"`
	Environments map[string]env `yaml:"environments"`
}

type env struct {
	Release string                 `yaml:"release"`
	Preview string                 `yaml:"preview"`
	Env     map[string]interface{} `yaml:"env"`
}

//
// Methods
//

// NewConfig creates new Config object
func NewConfig() *Config {
	return &Config{
		Wd: Getwd,
	}
}

// Env returns environment by name that have been described in config
// if no environment perstits in config.yaml it will return error
func (c Config) Env(name string) (env, error) {
	e, ok := c.YML.Environments[name]
	if !ok {
		return env{}, fmt.Errorf("Unknown environment: %s", name)
	}

	return e, nil
}

// Parse parsing config file
func (c *Config) Parse(file string) error {
	if file == "" || !path.IsAbs(file) {
		cwd, err := c.Wd()
		if err != nil {
			return err
		}

		name := "scruffy.yml"
		if file != "" {
			name = file
		}
		file = path.Join(cwd, name)
	}

	yml, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	var config configYML
	err = yaml.Unmarshal(yml, &config)
	if err != nil {
		return err
	}

	c.YML = &config

	return nil
}
