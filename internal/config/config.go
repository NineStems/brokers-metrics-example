package config

import (
	"gopkg.in/yaml.v3"
	"os"

	"mb-and-metrics/pkg/errors" //nolint
)

// Logger конфигурация для логгера.
type Logger struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
}

type Rest struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Server конфигурация для HTTP сервера.
type Server struct {
	Http Rest `yaml:"rest"`
}

// Kafka конфигурация очереди.
type Kafka struct {
	Brokers string `yaml:"brokers"`
	Topic   string `yaml:"topic"`
	Main    string `yaml:"main"`
}

// Rabbit конфигурация очереди.
type Rabbit struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Exchange   string `yaml:"exchange"`
	Queue      string `yaml:"queue"`
	Key        string `yaml:"key"`
	Credential struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"credential"`
	Main string `yaml:"main"`
}

// Config конфигурация сервиса.
type Config struct {
	Logger Logger `yaml:"logger"`
	Kafka  Kafka  `yaml:"kafka"`
	Rabbit Rabbit `yaml:"rabbit"`
	Server Server `yaml:"server"`
}

// Apply применяет значение из конфигурационного файла.
func (c *Config) Apply(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "os.Open")
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f) //nolint: typecheck,nolintlint
	if err = decoder.Decode(c); err != nil {
		return errors.Wrap(err, "decoder.Decode")
	}

	return nil
}

func New() *Config {
	return &Config{}
}
