package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Glowman554/infrastructure/config"
	"github.com/Glowman554/infrastructure/service"
	"github.com/Glowman554/infrastructure/utils"
)

type Subcommand func(config *config.Config) error

var Subcommands = map[string]Subcommand{
	"stop": func(config *config.Config) error {
		err := service.RunForServices(config, true, func(n string, s *service.Service) error {
			if s.Containers == nil {
				return nil
			}

			for _, c := range s.Containers {
				service.DeleteContainer(c)
			}
			return nil
		})
		if err != nil {
			return err
		}

		err = service.RunForServices(config, true, func(n string, s *service.Service) error {
			if s.Networks == nil {
				return nil
			}

			for _, n := range s.Networks {
				service.DeleteNetwork(n)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	},

	"start": func(config *config.Config) error {
		err := service.RunForServices(config, false, func(n string, s *service.Service) error {
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
		})
		if err != nil {
			return err
		}

		err = service.RunForServices(config, false, func(n string, s *service.Service) error {
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
		})
		if err != nil {
			return err
		}

		return nil
	},

	"build": func(config *config.Config) error {
		err := service.RunForServices(config, false, func(n string, s *service.Service) error {
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

				utils.Execute(*command, *directory, true)
			}

			return nil
		})
		if err != nil {
			return err
		}
		return nil
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
			err = cmd(config)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("Subcommand " + os.Args[i] + " not found")
		}
	}

}
