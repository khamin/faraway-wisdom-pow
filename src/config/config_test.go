package config

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestApply(t *testing.T) {
	Apply()

	if v := viper.GetString("log_fmt"); v != "text" {
		t.Errorf("default log format is %s", v)
	}

	if v := viper.GetString("log_level"); v != "info" {
		t.Errorf("default log level is %s", v)
	}
}

func TestSetLogFmt(t *testing.T) {
	std := logrus.StandardLogger()

	for format, formatter := range logFormats {
		viper.Set("log_fmt", strings.ToUpper(format))

		if err := setLogFmt(); err != nil {
			t.Error(err)
		}

		if std.Formatter != formatter {
			t.Error(std.Formatter)
		}
	}
}
