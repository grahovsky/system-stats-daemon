package config

import (
	"github.com/grahovsky/system-stats-daemon/internal/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Log struct {
		Level string `mapstructure:"level" env:"LOG_LEVEL"`
	} `mapstructure:"log"`
	Server struct {
		Host string `mapstructure:"host" env:"SRV_HOST"`
		Port string `mapstructure:"port" env:"HTTP_PORT"`
	} `mapstructure:"server"`
	Stats struct {
		Limit   int64 `mapstructure:"limit" env:"STATS_LIMIT"`
		Collect struct {
			LoadAverage bool `mapstructure:"loadAverage" env:"LOAD_AVERAGE"`
			CPU         bool `mapstructure:"cpu" env:"CPU"`
			DiskInfo    bool `mapstructure:"discInfo" env:"DISC_INFO"`
		} `mapstructure:"collect"`
	} `mapstructure:"stats"`
}

var Settings *Config

func init() {
	defaultSettings := defaultSettings()
	Settings = &defaultSettings

	pflag.String("loglevel", "INFO", "log level app")
	pflag.String("config", "./configs/config.yaml", "Path to configuration file")
	pflag.String("server_host", "0.0.0.0", "server hostname")
	pflag.String("server_port", "8086", "server port")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		logger.Error(err.Error())
	}

	viper.SetConfigFile(viper.Get("config").(string))
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	if err := viper.ReadInConfig(); err != nil {
		logger.Error(err.Error())
	}

	if err := viper.Unmarshal(&Settings); err != nil {
		logger.Error(err.Error())
	}

	envLogLevel := viper.Get("LOG_LEVEL")
	if envLogLevel != nil {
		Settings.Log.Level = envLogLevel.(string)
	}
}

func defaultSettings() Config {
	return Config{
		Log: struct {
			Level string `mapstructure:"level" env:"LOG_LEVEL"`
		}{Level: "DEBUG"},
		Server: struct {
			Host string `mapstructure:"host" env:"SRV_HOST"`
			Port string `mapstructure:"port" env:"HTTP_PORT"`
		}{Host: "0.0.0.0", Port: "8086"},
		Stats: struct {
			Limit   int64 "mapstructure:\"limit\" env:\"STATS_LIMIT\""
			Collect struct {
				LoadAverage bool "mapstructure:\"loadAverage\" env:\"LOAD_AVERAGE\""
				CPU         bool "mapstructure:\"cpu\" env:\"CPU\""
				DiskInfo    bool "mapstructure:\"discInfo\" env:\"DISC_INFO\""
			} "mapstructure:\"collect\""
		}{Limit: 600, Collect: struct {
			LoadAverage bool "mapstructure:\"loadAverage\" env:\"LOAD_AVERAGE\""
			CPU         bool "mapstructure:\"cpu\" env:\"CPU\""
			DiskInfo    bool "mapstructure:\"discInfo\" env:\"DISC_INFO\""
		}{LoadAverage: true, CPU: true, DiskInfo: true}},
	}
}
