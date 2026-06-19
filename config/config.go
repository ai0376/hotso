package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

//Mongo config
type Mongo struct {
	Host string `json:"host"`
}

type WebDav struct {
	Host      string   `json:"host"`
	User      string   `json:"user"`
	Password  string   `json:"password"`
	Path      string   `json:"data_path"`
	Files     []string `json:"files"`
	RemoteDir string   `json:"remote_dir"`
}

type Service struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type Redis struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Auth string `json:"password"`
}

type HotTop struct {
	BeginTime       int64 `json:"begin_unix"`
	DurationTimeSec int64 `json:"duration_sec"`
}

type Config struct {
	MongoDB Mongo   `json:"mongodb"`
	WebDav  WebDav  `json:"webdav"`
	Service Service `json:"service"`
	Redis   Redis   `json:"redis"`
	HotTop  HotTop  `json:"hottop"`
}

var (
	cfg     *Config
	loadErr error
	once    sync.Once
)

const defaultConfigFile = "../config/config.json"

func getCfgFile() string {
	if envPath := os.Getenv("HOTSO_CONFIG_PATH"); envPath != "" {
		return envPath
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return filepath.Join(dir, defaultConfigFile)
}

//LoadConfig loads config and returns error instead of panic
func LoadConfig() (*Config, error) {
	once.Do(func() {
		cfg = &Config{}
		path := getCfgFile()
		if path == "" {
			loadErr = fmt.Errorf("failed to determine config path")
			return
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			loadErr = err
			return
		}
		loadErr = json.Unmarshal(data, cfg)
	})
	if loadErr != nil {
		return nil, loadErr
	}
	return cfg, nil
}

//GetConfig get config
func GetConfig() *Config {
	c, err := LoadConfig()
	if err != nil {
		panic(err.Error())
	}
	return c
}
