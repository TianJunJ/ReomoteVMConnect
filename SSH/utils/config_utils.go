package utils

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

// Config  定义配置对象
var Config struct {
	Server struct {
		ServerPort string `yaml:"server-port"`
	} `yaml:"server"`
}

func InitConfig(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &Config)
	if err != nil {
		return err
	}
	return nil
}
