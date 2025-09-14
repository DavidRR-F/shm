package config

type Config struct {
	Links    []SymbolicLink `yaml:"links"`
	Managers []string       `yaml:"managers"`
}
