package config

import (
	"flag"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is a application configuration structure
type Config struct {
	ConfigConfluence struct {
        URL      string `yaml:"url" env:"CONFLUENCE_URL"`
        Username string `yaml:"username" env:"CONFLUENCE_USERNAME"`
        Token    string `yaml:"token" env:"CONFLUENCE_TOKEN"`
    }`yaml:"confluence"`
    ConfigSlack struct {
        Token    string `yaml:"token" env:"SLACK_TOKEN"`
    }`yaml:"slack"`
	Server struct {
        ActionURL string `yaml:"actionurl" env:"SERVER_ACTION_URL,ACTION_URL" env-description:"The action callback url for slack buttons." env-default:"/actions"`
		Port int `yaml:"port" env:"SERVER_PORT,PORT" env-description:"Server port" env-default:"3000"`
	} `yaml:"server"`
	Debug bool `yaml:"debug" env:"SERVER_DEBUG,DEBUG" env-description:"Enable debug output" env-default:"false"`
}

var (
    configFile = flag.String("config-file", "config.yaml", "The name and location of the configuration file.")
)

func Configuration() (*Config, error) {
    flag.Parse()
    var cfg Config
    err := cleanenv.ReadConfig(*configFile, &cfg)
    if err != nil {
        fmt.Printf("Could not read configuration: %v", err)
        return nil, err
    }
    return &cfg, nil
}
