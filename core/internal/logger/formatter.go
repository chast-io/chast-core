package chastlog

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
)

// Based on https://github.com/t-tomalak/logrus-easy-formatter and https://github.com/x-cray/logrus-prefixed-formatter

const (
	defaultLogFormat       = "%time% [%lvl%]: %caller% - %msg%\n"
	defaultTimestampFormat = "2006-01-02 15:04:05" // alternative use time.RFC3339
)

const (
	LEVEL  = "%lvl%"
	TIME   = "%time%"
	MSG    = "%msg%"
	CALLER = "%caller%"
	FILE   = "%file%"
)

var (
	defaultColorScheme = &ColorScheme{ //nolint:gochecknoglobals // only used internally
		PanicLevelStyle: "red",
		FatalLevelStyle: "red",
		ErrorLevelStyle: "red",
		WarnLevelStyle:  "yellow",
		DebugLevelStyle: "blue",
		InfoLevelStyle:  "green",
		TraceLevelStyle: "gray",

		TimestampStyle: "black+h",
	}
)

type ColorScheme struct {
	PanicLevelStyle string
	FatalLevelStyle string
	ErrorLevelStyle string
	WarnLevelStyle  string
	InfoLevelStyle  string
	DebugLevelStyle string
	TraceLevelStyle string

	TimestampStyle string
}

type compiledColorScheme struct {
	PanicLevelColor func(string) string
	FatalLevelColor func(string) string
	ErrorLevelColor func(string) string
	WarnLevelColor  func(string) string
	InfoLevelColor  func(string) string
	DebugLevelColor func(string) string
	TraceLevelColor func(string) string

	TimestampColor func(string) string
}

func getCompiledColor(main string, fallback string) func(string) string {
	var style string
	if main != "" {
		style = main
	} else {
		style = fallback
	}

	return ansi.ColorFunc(style)
}

func compileColorScheme(scheme *ColorScheme) *compiledColorScheme {
	return &compiledColorScheme{
		PanicLevelColor: getCompiledColor(scheme.PanicLevelStyle, defaultColorScheme.PanicLevelStyle),
		FatalLevelColor: getCompiledColor(scheme.FatalLevelStyle, defaultColorScheme.FatalLevelStyle),
		ErrorLevelColor: getCompiledColor(scheme.ErrorLevelStyle, defaultColorScheme.ErrorLevelStyle),
		WarnLevelColor:  getCompiledColor(scheme.WarnLevelStyle, defaultColorScheme.WarnLevelStyle),
		InfoLevelColor:  getCompiledColor(scheme.InfoLevelStyle, defaultColorScheme.InfoLevelStyle),
		DebugLevelColor: getCompiledColor(scheme.DebugLevelStyle, defaultColorScheme.DebugLevelStyle),
		TraceLevelColor: getCompiledColor(scheme.TraceLevelStyle, defaultColorScheme.TraceLevelStyle),

		TimestampColor: getCompiledColor(scheme.TimestampStyle, defaultColorScheme.TimestampStyle),
	}
}

// Formatter implements logrus.Formatter interface.
type TextFormatter struct {
	colorScheme *compiledColorScheme
}

func (f *TextFormatter) SetColorScheme(colorScheme *ColorScheme) {
	f.colorScheme = compileColorScheme(colorScheme)
}

func (f *TextFormatter) getColorScheme() *compiledColorScheme {
	if f.colorScheme == nil {
		f.colorScheme = compileColorScheme(defaultColorScheme)
	}

	return f.colorScheme
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := defaultLogFormat
	output = f.setTimestamp(output, entry)
	output = f.SetLevel(output, entry)

	output = strings.Replace(output, MSG, entry.Message, 1)

	output = strings.Replace(output, CALLER, filepath.Base(entry.Caller.Function), 1)

	output = strings.Replace(output, FILE, entry.Caller.File+":"+strconv.Itoa(entry.Caller.Line), 1)

	for key, val := range entry.Data {
		switch value := val.(type) {
		case string:
			output = strings.Replace(output, "%"+key+"%", value, 1)
		case int:
			s := strconv.Itoa(value)
			output = strings.Replace(output, "%"+key+"%", s, 1)
		case bool:
			s := strconv.FormatBool(value)
			output = strings.Replace(output, "%"+key+"%", s, 1)
		}
	}

	return []byte(output), nil
}

func (f *TextFormatter) setTimestamp(output string, entry *logrus.Entry) string {
	timestampFormat := defaultTimestampFormat
	formattedTimestamp := f.getColorScheme().TimestampColor(entry.Time.Format(timestampFormat))

	return strings.Replace(output, TIME, formattedTimestamp, 1)
}

func (f *TextFormatter) SetLevel(output string, entry *logrus.Entry) string {
	var color func(string) string

	level := strings.ToUpper(entry.Level.String())

	switch entry.Level {
	case logrus.PanicLevel:
		color = f.colorScheme.PanicLevelColor
	case logrus.FatalLevel:
		color = f.colorScheme.FatalLevelColor
	case logrus.ErrorLevel:
		color = f.colorScheme.ErrorLevelColor
	case logrus.WarnLevel:
		color = f.colorScheme.WarnLevelColor
	case logrus.InfoLevel:
		color = f.colorScheme.InfoLevelColor
	case logrus.DebugLevel:
		color = f.colorScheme.DebugLevelColor
	case logrus.TraceLevel:
		color = f.colorScheme.TraceLevelColor
	}

	level = color(level)

	if entry.Level < logrus.InfoLevel {
		output = strings.Replace(output, LEVEL, level, 1)
		output = strings.Replace(output, MSG, color(entry.Message), 1)
	} else {
		output = strings.Replace(output, LEVEL, level, 1)
	}

	return output
}
