package chastlog

import (
	"github.com/sirupsen/logrus"
)

var (
	Log = logrus.New() //nolint:gochecknoglobals // Logger abstraction
)

const (
	PanicLevel logrus.Level = logrus.PanicLevel
	FatalLevel              = logrus.FatalLevel
	ErrorLevel              = logrus.ErrorLevel
	WarnLevel               = logrus.WarnLevel
	InfoLevel               = logrus.InfoLevel
	DebugLevel              = logrus.DebugLevel
	TraceLevel              = logrus.TraceLevel
)

//nolint:gochecknoinits // Preconfigured logging settings
func init() {
	Log.SetReportCaller(true)

	formatter := TextFormatter{} //nolint:exhaustruct // using defaults
	Log.SetFormatter(&formatter)

	Log.SetLevel(TraceLevel)
}
