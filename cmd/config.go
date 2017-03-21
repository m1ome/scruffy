package cmd

import (
	"fmt"
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

type Config struct {
	YML *ConfigYML
	Wd  WorkingDirGetter
}

type ConfigYML struct {
	Token        string         `yaml:"token"`
	Source       string         `yaml:"source"`
	Environments map[string]Env `yaml:"environments"`
}

type Env struct {
	Release string                 `yaml:"release"`
	Preview string                 `yaml:"preview"`
	Env     map[string]interface{} `yaml:"env"`
}

//
// Methods
//
func NewConfig() *Config {
	return &Config{
		Wd: Getwd,
	}
}

func (c Config) Env(name string) (Env, error) {
	env, ok := c.YML.Environments[name]
	if !ok {
		return Env{}, fmt.Errorf("Unknown environment: %s", name)
	}

	return env, nil
}

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

	var config ConfigYML
	err = yaml.Unmarshal(yml, &config)
	if err != nil {
		return err
	}

	c.YML = &config

	return nil
}
