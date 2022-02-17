package config

import (
	"flag"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Apply() {
	flag.String("log_fmt", "text", "Log format")
	flag.String("log_level", "info", "Log level")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()

	logrus.SetReportCaller(true)

	if err := setLogFmt(); err != nil {
		logrus.Error(err)
	}

	if err := setLogLevel(); err != nil {
		logrus.Error(err)
	}
}
