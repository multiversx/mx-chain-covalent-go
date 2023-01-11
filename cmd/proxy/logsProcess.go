package main

import (
	"fmt"
	"os"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-covalent-go/cmd/proxy/config"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-logger-go/file"
)

const (
	defaultLogsPath      = "logs"
	logFileLifeSpanInSec = 86400
	logFileMaxSizeInMB   = 1024
)

func attachFileLogger(log logger.Logger, logFilePrefix string, flagsConfig config.FlagsLog) error {
	var err error
	if flagsConfig.SaveLogFile {
		fileLogging, err := file.NewFileLogging(file.ArgsFileLogging{
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
