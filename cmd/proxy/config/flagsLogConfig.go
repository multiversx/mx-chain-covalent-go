package config

// FlagsLog holds the log flags
type FlagsLog struct {
	WorkingDir       string
	LogLevel         string
	DisableAnsiColor bool
	SaveLogFile      bool
	EnableLogName    bool
}
