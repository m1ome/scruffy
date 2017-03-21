package cmd

import (
	"errors"
	"fmt"
	"github.com/m1ome/apiary"
)

type Publisher struct {
	Wd     WorkingDirGetter
	Parser *Parser
	Apiary apiary.ApiaryInterface
}

func NewPublisher(token string) *Publisher {
	return &Publisher{
		Wd:     Getwd,
		Parser: NewParser(),
		Apiary: apiary.NewApiary(apiary.ApiaryOptions{
			Token: token,
		}),
	}
}

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
