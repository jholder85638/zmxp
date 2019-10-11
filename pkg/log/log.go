package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jholder85638/zmxp/pkg/config"
	"github.com/jesseduffield/rollrus"
	"github.com/sirupsen/logrus"
)

// NewLogger returns a new logger
func NewLogger(config *config.AppConfig, rollrusHook string) *logrus.Entry {
	var log *logrus.Logger
	environment := "production"
	if true || config.Debug || os.Getenv("DEBUG") == "TRUE" { // TODO: remove true here
		environment = "development"
		log = newDevelopmentLogger(config)
	} else {
		log = newProductionLogger()
	}

	// highly recommended: tail -f development.log | humanlog
	// https://github.com/aybabtme/humanlog
	log.Formatter = &logrus.JSONFormatter{}

	if config.UserConfig.Reporting == "on" && rollrusHook != "" {
		// this isn't really a secret token: it only has permission to push new rollbar items
		hook := rollrus.NewHook(rollrusHook, environment)
		log.Hooks.Add(hook)
	}
	return log.WithFields(logrus.Fields{
		"debug":     config.Debug,
		"version":   config.Version,
		"commit":    config.Commit,
		"buildDate": config.BuildDate,
	})
}

func getLogLevel() logrus.Level {
	strLevel := os.Getenv("LOG_LEVEL")
	level, err := logrus.ParseLevel(strLevel)
	if err != nil {
		return logrus.DebugLevel
	}
	return level
}

func newDevelopmentLogger(config *config.AppConfig) *logrus.Logger {
	log := logrus.New()
	log.SetLevel(getLogLevel())
	file, err := os.OpenFile(filepath.Join(config.ConfigDir, "development.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("unable to log to file")
		os.Exit(1)
	}
	log.SetOutput(file)
	return log
}

func newProductionLogger() *logrus.Logger {
	log := logrus.New()
	log.Out = ioutil.Discard
	log.SetLevel(logrus.ErrorLevel)
	return log
}
