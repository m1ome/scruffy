package cmd

import (
	"errors"
	"testing"

	"github.com/m1ome/apiary"
	"strings"
	"os"
	"path"
)

type FakeApiary struct{}

func (a FakeApiary) Me() (me apiary.ApiaryMeResponse, err error) {
	return apiary.ApiaryMeResponse{}, nil
}
func (a FakeApiary) GetApis() (apis *apiary.ApiaryApisResponse, err error) {
	return nil, nil
}
func (a FakeApiary) GetTeamApis(team string) (apis *apiary.ApiaryApisResponse, err error) {
	return nil, nil
}
func (a FakeApiary) PublishBlueprint(name string, content []byte) (published bool, err error) {
	return false, errors.New("APIARY_ERROR")
}
func (a FakeApiary) FetchBlueprint(name string) (blueprint *apiary.ApiaryFetchResponse, err error) {
	return nil, nil
}

type FakeApiaryPublish struct {
	FakeApiary
}

func (a FakeApiaryPublish) PublishBlueprint(name string, content []byte) (published bool, err error) {
	return true, nil
}

type FakeApiaryNonPublish struct {
	FakeApiary
}

func (a FakeApiaryNonPublish) PublishBlueprint(name string, content []byte) (published bool, err error) {
	return false, nil
}

func TestPublish(t *testing.T) {
	t.Run("Parsing error", func(t *testing.T) {
		p := NewPublisher("token")
		err := p.Publish("/unknown/directory", "wrong_name", nil)

		if err == nil {
			t.Error("Wrong directory should return error")
		}
	})

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	sources := path.Join(cwd, "test", "input", "version_1")

	t.Run("Apiary error", func(t *testing.T) {
		p := &Publisher{
			Wd:     Getwd,
			Parser: NewParser(),
			Apiary: &FakeApiary{},
		}

		config := NewConfig()
		config.Parse("test/Config.yml")
		env, err := config.Env("public")
		if err != nil {
			t.Fatal(err)
		}

		err = p.Publish(sources, env.Release, env.Env)
		if err == nil || !strings.Contains(err.Error(), "APIARY_ERROR") {
			t.Errorf("Not return error: %s", err)
		}
	})

	t.Run("Apiary publish", func(t *testing.T) {
		p := &Publisher{
			Wd:     Getwd,
			Parser: NewParser(),
			Apiary: &FakeApiaryPublish{},
		}

		config := NewConfig()
		config.Parse("test/Config.yml")
		env, err := config.Env("public")
		if err != nil {
			t.Fatal(err)
		}

		err = p.Publish(sources, env.Release, env.Env)
		if err != nil {
			t.Errorf("Returns error: %s", err)
		}
	})

	t.Run("Apiary publish", func(t *testing.T) {
		p := &Publisher{
			Wd:     Getwd,
			Parser: NewParser(),
			Apiary: &FakeApiaryNonPublish{},
		}

		config := NewConfig()
		config.Parse("test/Config.yml")
		env, err := config.Env("public")
		if err != nil {
			t.Fatal(err)
		}

		err = p.Publish(sources, env.Release, env.Env)
		if err == nil {
			t.Error("Should return publish error")
		}
	})
}
