package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"common/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "supervisor",
	Short: "Workspace container init process",
}

var (
	// ServiceName is the name we use for tracing/logging.
	ServiceName = "supervisor"
	// Version of this service - set during build.
	Version = ""
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	log.Init(ServiceName, Version, true, false)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type fatalTerminationLogHook struct{}

func (fatalTerminationLogHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.FatalLevel}
}

func (fatalTerminationLogHook) Fire(e *logrus.Entry) error {
	msg := e.Message
	if err := e.Data[logrus.ErrorKey]; err != nil {
		msg += ": " + err.(error).Error()
	}

	return os.WriteFile("/tmp/opencoder/termination", []byte(msg), 0o644)
}

const (
	opencoderLogDir       = "/tmp/opencoder/log"
	supervisorLogFilePath = opencoderLogDir + "/supervisor.log"
)

func initLog(json bool) io.Closer {
	log.Init(ServiceName, Version, json, os.Getenv("SUPERVISOR_DEBUG_ENABLE") == "true")
	log.Log.Logger.AddHook(fatalTerminationLogHook{})
	if err := os.MkdirAll(opencoderLogDir, 0755); err != nil {
		log.WithError(err).Error("cannot create opencoder log directory")
		return io.NopCloser(nil)
	}
	supervisorLogFile, err := os.OpenFile(supervisorLogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.WithError(err).Error("cannot open supervisor log file")
		return io.NopCloser(nil)
	}
	log.Log.Logger.AddHook(&FileHook{
		Writer: supervisorLogFile,
		Formatter: &prefixed.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			FullTimestamp:   true,
		},
	})
	return supervisorLogFile
}

type FileHook struct {
	Writer    io.Writer
	Formatter logrus.Formatter
}

func (hook *FileHook) Fire(entry *logrus.Entry) error {
	formatted, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(formatted)
	return err
}

func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
