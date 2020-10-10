package config

import "time"

type Database struct {
	Kind           string        `yaml:"kind"`
	Host           string        `yaml:"host"`
	Port           string        `yaml:"port"`
	User           string        `yaml:"user"`
	Pass           string        `yaml:"pass"`
	DBName         string        `yaml:"dbname"`
	MaxConnections int           `yaml:"max_connections"`
	Timeout        time.Duration `yaml:"timeout"`
}
