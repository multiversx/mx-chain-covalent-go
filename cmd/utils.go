package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ElrondNetwork/covalent-indexer-go/cmd/config"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/common/logging"
)

const (
	defaultLogsPath      = "logs"
	logFileLifeSpanInSec = 86400
	logFileMaxSizeInMB   = 1024
)

func attachFileLogger(log logger.Logger, logFilePrefix string, flagsConfig config.FLagsLog) error {
	var err error
	if flagsConfig.SaveLogFile {
		fileLogging, err := logging.NewFileLogging(logging.ArgsFileLogging{
			WorkingDir:      flagsConfig.WorkingDir,
			DefaultLogsPath: defaultLogsPath,
			LogFilePrefix:   logFilePrefix,
		})
		if err != nil {
			return fmt.Errorf("%w creating a log file", err)
		}

		if !check.IfNil(fileLogging) {
			err = fileLogging.ChangeFileLifeSpan(time.Second*time.Duration(logFileLifeSpanInSec), logFileMaxSizeInMB)
			if err != nil {
				return err
			}
		}
	}

	err = logger.SetDisplayByteSlice(logger.ToHex)
	log.LogIfError(err)
	logger.ToggleLoggerName(flagsConfig.EnableLogName)
	logLevelFlagValue := flagsConfig.LogLevel
	err = logger.SetLogLevel(logLevelFlagValue)
	if err != nil {
		return err
	}

	if flagsConfig.DisableAnsiColor {
		err = logger.RemoveLogObserver(os.Stdout)
		if err != nil {
			return err
		}

		err = logger.AddLogObserver(os.Stdout, &logger.PlainFormatter{})
		if err != nil {
			return err
		}
	}
	log.Trace("logger updated", "level", logLevelFlagValue, "disable ANSI color", flagsConfig.DisableAnsiColor)

	return nil
}
