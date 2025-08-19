package service

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/Glowman554/infrastructure/utils"
)

func CreateNetwork(network Network) error {
	networkType := "bridge"
	if network.Type != nil {
		networkType = *network.Type
	}

	slog.Info("creating network "+network.Name, "type", networkType)

	return utils.Execute("sudo docker network create "+network.Name+" --driver "+networkType, "/", false)
}

func DeleteNetwork(network Network) error {
	slog.Info("removing network " + network.Name)
	return utils.Execute("sudo docker network rm "+network.Name, "/", false)
}

func CreateContainer(container Container, service string, secrets map[string]string) error {
	slog.Info("creating container " + container.Name)

	networks := []string{}
	if container.Networks != nil {
		for _, network := range container.Networks {
			networks = append(networks, "--network "+network)
		}
	}

	environment := []string{}
	if container.Environment != nil {
		for key, value := range container.Environment {
			environment = append(environment, "--env "+key+"="+value)
		}
	}

	ports := []string{}
	if container.Ports != nil {
		for host, guest := range container.Ports {
			ports = append(ports, "--publish "+host+":"+guest)
		}
	}

	volumes := []string{}
	if container.Mounts != nil {
		for host, guest := range container.Mounts {
			volumes = append(volumes, "--volume "+host+":"+guest)
		}
	}

	aliases := []string{}
	if container.Aliases != nil {
		for _, alias := range container.Aliases {
			aliases = append(aliases, "--network-alias "+alias)
		}
	}

	privileged := ""
	if container.Privileged != nil && *container.Privileged {
		privileged = "--privileged"
	}

	command := ""
	if container.Command != nil {
		command = *container.Command
	}

	userID := 0
	if container.UserID != nil {
		userID = *container.UserID
	}

	dockerCommand := "sudo docker run -d -it --restart always " +
		strings.Join(networks, " ") + " " +
		strings.Join(environment, " ") + " " +
		strings.Join(ports, " ") + " " +
		strings.Join(volumes, " ") + " " +
		strings.Join(aliases, " ") + " " +
		privileged + " " +
		"--name " + container.Name + " --hostname " + container.Name + " " +
		"--user " + strconv.Itoa(userID) + ":" + strconv.Itoa(userID) + " " +
		container.Image + " " + command

	tmp, err := ReplaceAll(dockerCommand, service, secrets)
	if err != nil {
		return err
	}

	return utils.Execute(*tmp, "/", true)
}

func DeleteContainer(container Container) error {
	slog.Info("stopping and removing container " + container.Name)

	err := utils.Execute("sudo docker stop "+container.Name, "/", false)
	if err != nil {
		return err
	}
	return utils.Execute("sudo docker rm "+container.Name, "/", false)
}
