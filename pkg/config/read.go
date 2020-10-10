package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func Read(path string, cfg interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Errorf("cant read config file: %s", err.Error())
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return errors.Errorf("cant parse config: %s", err.Error())
	}

	return nil
}
