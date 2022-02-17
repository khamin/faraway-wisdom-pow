package config

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var logFormats = map[string]logrus.Formatter{
	"json": &logrus.JSONFormatter{},
	"text": &logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	},
}

func setLogFmt() (err error) {
	format := viper.GetString("log_fmt")
	format = strings.ToLower(format)
	formatter, ok := logFormats[format]

	if !ok {
		return fmt.Errorf("unknown log format %s", format)
	}

	logrus.SetFormatter(formatter)
	return
}

func setLogLevel() (err error) {
	level := viper.GetString("log_level")
	id, err := logrus.ParseLevel(level)

	if err == nil {
		logrus.SetLevel(id)
	}

	return
}
