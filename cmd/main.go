package main

import (
	"github.com/sirupsen/logrus"
	"github.com/zakimal/minesweeper/apiserver"
)

func main() {
	logger := logrus.StandardLogger()
	logger.Infoln("Starting API Server...")
	if err := apiserver.Start(logger); err != nil {
		logger.WithError(err).Fatalln("can not start API server")
	}
}
