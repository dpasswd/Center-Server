package config

import "github.com/spf13/viper"

type Logger struct {
	Path       string
	Level      string
	Stdout     bool
	EnabledWeb bool
	EnabledDB  bool
}

func InitLog(cfg *viper.Viper) *Logger {
	return &Logger{
		Path:       cfg.GetString("path"),
		Level:      cfg.GetString("level"),
		Stdout:     cfg.GetBool("stdout"),
		EnabledWeb: cfg.GetBool("enabledWeb"),
		EnabledDB:  cfg.GetBool("enabledDb"),
	}
}

var LoggerConfig = new(Logger)
