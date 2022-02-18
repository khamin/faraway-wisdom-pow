package main

import (
	"flag"
	"sync"
	"wisdom/client/request"
	"wisdom/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	flag.String("server", "127.0.0.1:8000", "Server address")
	flag.Uint("limit", 1, "Requests limit")

	config.Apply()

	server := viper.GetString("server")
	limit := viper.GetInt("limit")
	logFields := logrus.WithField("limit", limit)

	var wg sync.WaitGroup
	wg.Add(limit)

	for i := 0; i < limit; i++ {
		go func() {
			defer wg.Done()

			quote, err := request.Run(server)

			if err != nil {
				panic(err)
			}

			logrus.WithField("quote", string(quote)).Info("quote received")
		}()
	}

	wg.Wait()
	logFields.Info("done")
}
