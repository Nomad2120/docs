package common

import (
	"os"
	"path"

	"gitlab.enterprise.qazafn.kz/oci/oci-docs/config"
	model "gitlab.enterprise.qazafn.kz/oci/oci-docs/model"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger -
func InitLogger(conf *config.Config) *zap.SugaredLogger {
	cfg := zap.Config{
		Encoding: "json",
		Level:    zap.NewAtomicLevelAt(getLogLevel(conf.LogLevel)),

		EncoderConfig: zapcore.EncoderConfig{
			//CallerKey: "caller",

			// FunctionKey: "func",

			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}

	if conf.LogDir != "" {
		err := os.MkdirAll(conf.LogDir, os.ModePerm)
		if err != nil {
			panic(errors.Wrapf(err, "logDir=%s", conf.LogDir))
		}
		cfg.OutputPaths = []string{path.Join(conf.LogDir, "out.log")}
		cfg.ErrorOutputPaths = []string{path.Join(conf.LogDir, "error.log")}
	} else {
		cfg.OutputPaths = []string{"stderr"}
		cfg.ErrorOutputPaths = []string{"stderr"}
	}

	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return log.Sugar().With("app", model.AppName)
}

func getLogLevel(text string) zapcore.Level {
	switch text {
	case "debug", "DEBUG":
		return zapcore.DebugLevel
	case "info", "INFO", "": // make the zero value useful
		return zapcore.InfoLevel
	case "warn", "WARN":
		return zapcore.WarnLevel
	case "error", "ERROR":
		return zapcore.ErrorLevel
	case "dpanic", "DPANIC":
		return zapcore.DPanicLevel
	case "panic", "PANIC":
		return zapcore.PanicLevel
	case "fatal", "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
