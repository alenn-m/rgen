package config

type Config struct {
	Package        string `yaml:"Package"`
	AutoMigrations bool   `yaml:"AutoMigrations"`
}
