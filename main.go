package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Glowman554/infrastructure/config"
	"github.com/Glowman554/infrastructure/service"
	"github.com/Glowman554/infrastructure/utils"
)

type SubcommandStep func(config *config.Config, n string, s *service.Service) error
type Subcommand struct {
	reverse bool
	steps   []SubcommandStep
}

var Subcommands = map[string]Subcommand{
	"stop": Subcommand{
		reverse: true,
		steps: []SubcommandStep{
			func(config *config.Config, n string, s *service.Service) error {
				if s.Containers == nil {
					return nil
				}

				for _, c := range s.Containers {
					service.DeleteContainer(c)
				}
				return nil
			},
			func(config *config.Config, n string, s *service.Service) error {
				if s.Networks == nil {
					return nil
				}

				for _, n := range s.Networks {
					service.DeleteNetwork(n)
				}
				return nil
			},
		},
	},

	"start": Subcommand{
		reverse: false,
		steps: []SubcommandStep{
			func(config *config.Config, n string, s *service.Service) error {
				if s.Networks == nil {
					return nil
				}

				for _, n := range s.Networks {
					err := service.CreateNetwork(n)
					if err != nil {
						return err
					}
				}
				return nil
			},
			func(config *config.Config, n string, s *service.Service) error {
				if s.Containers == nil {
					return nil
				}

				for _, c := range s.Containers {
					err := service.CreateContainer(c, n, config.Secrets)
					if err != nil {
						return err
					}
				}
				return nil
			},
		},
	},

	"build": Subcommand{
		reverse: false,
		steps: []SubcommandStep{
			func(config *config.Config, n string, s *service.Service) error {
				if s.Build == nil {
					return nil
				}

				for _, b := range s.Build {
					command, err := service.ReplaceAll(b.Command, n, config.Secrets)
					if err != nil {
						return err
					}

					directory, err := service.ReplaceAll(b.Directory, n, config.Secrets)
					if err != nil {
						return err
					}

					err = utils.Execute(*command, *directory, true)
					if err != nil {
						return err
					}
				}

				return nil
			},
		},
	},

	"clean": Subcommand{
		reverse: false,
		steps: []SubcommandStep{
			func(config *config.Config, n string, s *service.Service) error {
				if s.Containers == nil {
					return nil
				}

				for _, c := range s.Containers {
					err := service.DeleteImage(c)
					if err != nil {
						return err
					}
				}

				return nil
			},
		},
	},
}

func main() {
	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	if _, ok := os.LookupEnv("DEBUG"); ok {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "help" {
			fmt.Println("Available commands:")
			for key := range Subcommands {
				fmt.Println("> " + key)
			}
			return
		}
		
		if cmd, ok := Subcommands[os.Args[i]]; ok {
			if single := os.Getenv("SINGLE"); single != "" {
				s, err := service.LoadService(single)
				if err != nil {
					panic(err)
				}

				for _, j := range cmd.steps {
					err = j(config, single, s)
					if err != nil {
						panic(err)
					}
				}
			} else {
				for _, j := range cmd.steps {
					err = service.RunForServices(config, cmd.reverse, func(name string, service *service.Service) error {
						return j(config, name, service)
					})
					if err != nil {
						panic(err)
					}
				}
			}
		} else {
			fmt.Println("Subcommand " + os.Args[i] + " not found")
		}
	}

}
