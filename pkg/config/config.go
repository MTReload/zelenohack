package config

type Config struct {
	Database   string `yaml:"database"`
	ListenAddr string `yaml:"listen_addr"`
}
