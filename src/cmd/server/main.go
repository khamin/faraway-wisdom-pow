package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os/signal"
	"syscall"
	"time"
	"wisdom/config"
	"wisdom/quotes"
	"wisdom/server"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	flag.Duration("stop_timeout", 5*time.Second, "Shutdown timeout")
	flag.String("addr", "0.0.0.0:8000", "Bind address")
	flag.String("quotes", "/data/quotes.json", "Quotes input file")

	config.Apply()

	now := time.Now().UnixNano()
	rand.Seed(now)

	addr := viper.GetString("addr")
	quotesFile := viper.GetString("quotes")
	logFields := logrus.WithField("addr", addr)
	quotes, err := quotes.New(quotesFile)

	if err != nil {
		logFields.WithFields(logrus.Fields{
			"file": quotesFile,
			"err":  err,
		}).Panic("quotes load failure")
	}

	logFields.WithFields(logrus.Fields{
		"count": len(*quotes),
		"file":  quotesFile,
	}).Debug("quotes loaded")

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	defer stop()

	stopTimeout := viper.GetDuration("stop_timeout")

	config := &server.Config{
		Addr:        addr,
		Quotes:      quotes,
		StopTimeout: stopTimeout,
	}

	srv, err := server.New(config)

	if err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
	logFields.Info("stopping gracefully")

	if err := srv.Stop(); err != nil {
		logFields.WithField("err", err).Warn("forced to stop")
		return
	}

	logFields.Info("server stopped")
}
