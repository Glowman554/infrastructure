package config

type Config struct {
	Services []string          `json:"services"`
	Secrets  map[string]string `json:"secrets"`
}
