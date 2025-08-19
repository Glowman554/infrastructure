package service

type Network struct {
	Type *string `json:"type,omitempty"`
	Name string  `json:"name"`
}

type Container struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Networks    []string          `json:"networks,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Ports       map[string]string `json:"ports,omitempty"`
	Mounts      map[string]string `json:"mounts,omitempty"`
	Aliases     []string          `json:"aliases,omitempty"`
	Command     *string           `json:"command,omitempty"`
	Privileged  *bool             `json:"privileged,omitempty"`
	UserID      *int              `json:"userID,omitempty"`
}

type Build struct {
	Command   string `json:"command"`
	Directory string `json:"directory"`
}

type Service struct {
	Build      []Build     `json:"build,omitempty"`
	Networks   []Network   `json:"networks,omitempty"`
	Containers []Container `json:"containers,omitempty"`
}
