package config

// FLagsLog holds the log flags
type FLagsLog struct {
	WorkingDir       string
	LogLevel         string
	DisableAnsiColor bool
	SaveLogFile      bool
	EnableLogName    bool
}
