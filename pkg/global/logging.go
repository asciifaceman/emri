package global

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logLevel zap.AtomicLevel
)

func InitLogging(level zapcore.Level) error {
	logLevel.SetLevel(level)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "tz"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	c := zap.Config{
		Level:             logLevel,
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    map[string]interface{}{"pid": os.Getpid()},
	}
	conf, err := c.Build()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(conf.Named(App))
	zap.S().Debugw("initialized global logger", "level", level)
	return nil
}

func SetLogLevel(level zapcore.Level) {
	logLevel.SetLevel(level)
	zap.S().Debugw("global logger level updated", "new", logLevel.String())
}

func init() {
	logLevel = zap.NewAtomicLevel()
}
