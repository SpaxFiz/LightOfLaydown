// @program: unjuanable
// @author: Fizzy
// @created: 2021-12-02

package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"runtime"
)

var GlobalConfig *GlobalConf

type GlobalConf struct {
	Server struct {
		Cache struct {
			CachePath           string `yaml:"cache_path"`
			CacheDurationSecond int    `yaml:"cache_duration_second"`
		} `yaml:"cache"`
	} `yaml:"server"`
}

func NewGlobalConf(confPath string) (*GlobalConf, error) {
	var conf *GlobalConf
	file, err := os.Open(confPath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func init() {
	_, filepath, _, _ := runtime.Caller(0)
	dirPath := path.Dir(filepath)
	confPath := path.Join(dirPath, "..", "config.yaml")
	if conf, err := NewGlobalConf(confPath); err != nil {
		panic(err)
	} else {
		GlobalConfig = conf
	}
}
