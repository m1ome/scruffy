package cmd

import (
	"errors"
	"fmt"
	"github.com/m1ome/apiary"
)

func Publish(source, name, token string, env interface{}) error {

	buf, err := Parse(source, env)
	if err != nil {
		return errors.New(fmt.Sprintf("Parsing error: %s", err.Error()))
	}

	api := apiary.NewApiary(apiary.ApiaryOptions{
		Token: token,
	})

	published, err := api.PublishBlueprint(name, buf)
	if err != nil {
		return errors.New(fmt.Sprintf("Publishing error: %s", err.Error()))
	}

	if published {
		return nil
	}

	return errors.New("Source cannot be published")
}
