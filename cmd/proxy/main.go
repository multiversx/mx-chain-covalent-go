package main

import (
	"os"

	"github.com/urfave/cli"
)

const (
	logFilePrefix = "covalent-proxy"
)

func main() {
	app := cli.NewApp()
	app.Name = "Covalent proxy indexer tool"
	app.Usage = "This is the entry point for covalent proxy indexer tool. It acts as a proxy to fetch hyperblocks from Elrond  It converts hyperblocks data and provides API calls for covalent in their desired format"
	app.Flags = getFlags()
	app.Authors = []cli.Author{
		{
			Name:  "The Elrond Team",
			Email: "contact@elrond.com",
		},
	}

	app.Action = startProxy
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
		return
	}
}

func startProxy(ctx *cli.Context) error {
	flagsConfig := getFlagsLogConfig(ctx)
	errLogger := attachFileLogger(log, logFilePrefix, flagsConfig)
	if errLogger != nil {
		return errLogger
	}

	log.Info("starting server")

	return nil
}
