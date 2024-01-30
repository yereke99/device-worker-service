package config

type Config struct {
	LocalUrl string `yaml:"URL"`
}

func NewConfig(fileName string) (*Config, error) {

	cfg := new(Config)

	cfg.LocalUrl = "http://127.0.0.1:5555"

	return cfg, nil
}
