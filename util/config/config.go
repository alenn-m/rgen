package config

// Config represents application configuration
type Config struct {
	Package   string    `yaml:"Package"`
	Migration Migration `yaml:"Migration"`
}

// Migration represents migration configuration
type Migration struct {
	Sequential bool `yaml:"Sequential"`
}
