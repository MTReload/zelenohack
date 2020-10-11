package config

type Config struct {
	Database   string `yaml:"database"`
	ListenAddr string `yaml:"listen_addr"`
	FrontAddr  string `yaml:"front_addr"`
	Index      string `yaml:"index"`
}
