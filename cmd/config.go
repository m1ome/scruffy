package cmd

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	Token   string `yaml:"token"`
	Source  string `yaml:"source"`
	Public  Env    `yaml:"public"`
	Private Env    `yaml:"private"`
}

type Env struct {
	Name    string                 `yaml:"name"`
	Preview string                 `yaml:"preview"`
	Env     map[string]interface{} `yaml:"env"`
}

func ParseConfig(file string) (*Config, error) {
	if file == "" || !path.IsAbs(file) {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		name := "scruffy.yml"
		if file != "" {
			name = file
		}
		file = path.Join(cwd, name)
	}

	yml, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(yml, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
