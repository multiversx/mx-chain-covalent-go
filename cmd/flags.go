package main

import (
	"github.com/ElrondNetwork/covalent-indexer-go/cmd/config"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/urfave/cli"
)

var (
	workingDirectory = cli.StringFlag{
		Name:  "working-directory",
		Usage: "This flag specifies the `directory` where the application will use the logs.",
		Value: "",
	}
	logLevel = cli.StringFlag{
		Name: "log-level",
		Usage: "This flag specifies the logger `level(s)`. It can contain multiple comma-separated value. For example" +
			", if set to *:INFO the logs for all packages will have the INFO level. However, if set to *:INFO,api:DEBUG" +
			" the logs for all packages will have the INFO level, excepting the api package which will receive a DEBUG" +
			" log level.",
		Value: "*:" + logger.LogDebug.String(),
	}
	saveLogFile = cli.BoolFlag{
		Name:  "log-save",
		Usage: "Boolean option for enabling log saving. If set, it will automatically save all the logs into a file.",
	}
	enableLogName = cli.BoolFlag{
		Name:  "log-logger-name",
		Usage: "Boolean option to enable logger name in the logs.",
	}
	disableAnsiColor = cli.BoolFlag{
		Name:  "disable-ansi-color",
		Usage: "Boolean option for disabling ANSI colors in the logging system.",
	}
)

func getFlags() []cli.Flag {
	return []cli.Flag{
		workingDirectory,
		logLevel,
		saveLogFile,
		enableLogName,
		disableAnsiColor,
	}
}

func getFlagsLogConfig(ctx *cli.Context) config.FLagsLog {
	return config.FLagsLog{
		WorkingDir:       ctx.GlobalString(workingDirectory.Name),
		LogLevel:         ctx.GlobalString(logLevel.Name),
		DisableAnsiColor: ctx.GlobalBool(disableAnsiColor.Name),
		SaveLogFile:      ctx.GlobalBool(saveLogFile.Name),
		EnableLogName:    ctx.GlobalBool(enableLogName.Name),
	}
}
