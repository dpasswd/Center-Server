package config

import "github.com/spf13/viper"

type Upload struct {
	Path string
}

func InitUpload(cfg *viper.Viper) *Upload {
	return &Upload{
		Path: cfg.GetString("path"),
	}
}

var UploadConfig = new(Upload)
