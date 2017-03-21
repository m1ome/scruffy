// +build !test
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/m1ome/scruffy/cmd"
	"github.com/radovskyb/watcher"
	cli "gopkg.in/urfave/cli.v1"
)

func watch(config *cmd.Config, c *cli.Context, wf func(config *cmd.Config, c *cli.Context) error) error {

	fmt.Printf("Start watching changes in: %s\n", config.YML.Source)

	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write, watcher.Create)

	done := make(chan bool)
	watcherError := make(chan error)
	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Printf("Source directory[%s] changed, republishing\n", event.Path)
				config := cmd.NewConfig()
				err := config.Parse(c.String("config"))
				if err != nil {
					watcherError <- err
					done <- true

					return
				}

				err = wf(config, c)
				if err != nil {
					watcherError <- err
					done <- true

					return
				}
			case err := <-w.Error:
				watcherError <- err
				done <- true

				return
			case <-w.Closed:

				return
			}
		}
	}()

	err := w.AddRecursive(config.YML.Source)
	if err != nil {
		return fmt.Errorf("Adding folders to watch error: %s\n", err.Error())
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		return fmt.Errorf("Watcher error: %s\n", err.Error())
	}

	<-done
	return nil
}

func publishChanges(config *cmd.Config, c *cli.Context) (err error) {
	p := cmd.NewPublisher(config.YML.Token)

	env, err := config.Env(c.String("env"))
	if err != nil {
		err = fmt.Errorf("Environment error: %s", err)
		return
	}

	var name string
	if c.Bool("release") {
		name = env.Release
	} else {
		name = env.Preview
	}

	err = p.Publish(config.YML.Source, name, env.Env)

	if err != nil {
		err = fmt.Errorf("Publish error: %s", err)
		return
	}

	fmt.Printf("Public docs changed: http://docs.%s.apiary.io/#\n", name)

	return
}

func buildChanges(config *cmd.Config, c *cli.Context) (err error) {
	src := config.YML.Source
	env, err := config.Env(c.String("env"))

	if err != nil {
		return err
	}

	cwd, err := cmd.Getwd()
	if err != nil {
		return fmt.Errorf("Cwd error: %s", cwd)
	}

	build := path.Join(cwd, "build", "apiary.apib")
	if _, err := os.Stat(build); os.IsNotExist(err) {
		os.MkdirAll(path.Dir(build), 0770)
	}

	parser := cmd.NewParser()
	buf, err := parser.Parse(src, env.Env)
	if err != nil {
		return fmt.Errorf("Parsing error: %s", err.Error())
	}

	// Check if build file exists and remove it
	_, err = os.Stat(build)

	// create file if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(build)
		if err != nil {
			return fmt.Errorf("Create file error: %s", err.Error())
		}

		defer file.Close()
	}

	err = ioutil.WriteFile(build, buf, 0644)
	if err != nil {
		return fmt.Errorf("Building error: %s", err.Error())
	}

	fmt.Printf("Build avaiable at: %s\n", build)
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "Scruffy"
	app.Usage = "build your blueprints from mess to order!"
	app.Version = cmd.Version

	// Defining flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "token",
			Usage: "apiary.io token",
		},
	}

	// Defining commands
	app.Commands = []cli.Command{
		{
			Name:  "publish",
			Usage: "Publish/Preview your public blueprint",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Usage: "application configuration in yaml `config.yml`",
				},

				cli.StringFlag{
					Name:  "env",
					Usage: "Environment that have been set in config",
				},

				cli.BoolFlag{
					Name:  "release",
					Usage: "Release changes in production doc `false`",
				},

				cli.BoolFlag{
					Name:  "watch",
					Usage: "Watch changes and reload",
				},
			},
			Action: func(c *cli.Context) error {
				config := cmd.NewConfig()
				err := config.Parse(c.String("config"))
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("Config parsing error: %s\n", err.Error()), 1)
				}

				err = publishChanges(config, c)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}

				if c.Bool("watch") {
					err := watch(config, c, func(config *cmd.Config, c *cli.Context) error {
						return publishChanges(config, c)
					})

					if err != nil {
						return cli.NewExitError(err, 1)
					}
				}

				return nil
			},
		},

		{
			Name:  "build",
			Usage: "Build your blueprint",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Usage: "application configuration in yaml",
				},

				cli.StringFlag{
					Name:  "env",
					Usage: "Environment that have been set in config",
				},

				cli.BoolFlag{
					Name:  "watch",
					Usage: "Watch changes and reload on file change `false`",
				},
			},
			Action: func(c *cli.Context) error {
				config := cmd.NewConfig()
				err := config.Parse(c.String("config"))
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("Config parsing error: %s\n", err.Error()), 1)
				}

				err = buildChanges(config, c)
				if err != nil {
					return cli.NewExitError(err, 1)
				}

				if c.Bool("watch") {
					err := watch(config, c, func(config *cmd.Config, c *cli.Context) error {
						return buildChanges(config, c)
					})

					if err != nil {
						return cli.NewExitError(err, 1)
					}
				}

				return nil
			},
		},
	}

	app.Run(os.Args)
}
