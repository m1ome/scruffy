package cmd

import (
	"errors"
	"fmt"
	"github.com/m1ome/apiary"
)

// Publisher object, contains following:
//
// Parser - Parser object
// Apiary - http://github.com/m1ome/apiary API object
type Publisher struct {
	Wd     WorkingDirGetter
	Parser *Parser
	Apiary apiary.ApiaryInterface
}

// NewPublisher returns Publisher object
func NewPublisher(token string) *Publisher {
	return &Publisher{
		Wd:     Getwd,
		Parser: NewParser(),
		Apiary: apiary.NewApiary(apiary.ApiaryOptions{
			Token: token,
		}),
	}
}

// Publish build & publish template to apiary.io
func (p Publisher) Publish(source, name string, env interface{}) error {
	buf, err := p.Parser.Parse(source, env)
	if err != nil {
		return errors.New(fmt.Sprintf("Parsing error: %s", err.Error()))
	}

	published, err := p.Apiary.PublishBlueprint(name, buf)
	if err != nil {
		return errors.New(fmt.Sprintf("Publishing error: %s", err.Error()))
	}

	if published {
		return nil
	}

	return errors.New("Source cannot be published")
}
