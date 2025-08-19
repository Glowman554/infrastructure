package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Glowman554/infrastructure/config"
	"github.com/Glowman554/infrastructure/utils"
)

func LoadService(service string) (*Service, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(home, "Projects", service, ".toxicfox.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var s Service
	err = json.Unmarshal(data, &s)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal service: %w", err)
	}

	return &s, nil
}

type Executor func(name string, service *Service) error

func RunForServices(config *config.Config, reverse bool, executor Executor) error {
	services := config.Services
	if reverse {
		services = utils.Reverse(services)
	}

	for _, i := range services {
		service, err := LoadService(i)
		if err != nil {
			return err
		}

		err = executor(i, service)
		if err != nil {
			return err
		}
	}

	return nil
}
