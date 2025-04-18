package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// LogrusFields adalah alias untuk logrus.Fields yang sesuai
type LogrusFields logrus.Fields

func SetupLogger() {
	Log = logrus.New()

	// Set format log (bisa JSON atau Text)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set output log (defaultnya os.Stdout)
	Log.SetOutput(os.Stdout)

	// Set level log (bisa Debug, Info, Warn, Error, Fatal, Panic)
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		Log.SetLevel(logrus.InfoLevel) // Default level adalah Info
	} else {
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			Log.Warnf("Gagal memparse level log '%s', menggunakan level default Info", logLevel)
			Log.SetLevel(logrus.InfoLevel)
		} else {
			Log.SetLevel(level)
		}
	}
}
