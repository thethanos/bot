package logger

import (
	"bot/internal/config"
	"bytes"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Printer func(string, ...zap.Field)

func (p Printer) Write(b []byte) (int, error) {
	if p != nil {
		p(string(bytes.TrimSpace(b)))
	}
	return len(b), nil
}

type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

type ReleaseLogger struct {
	*zap.SugaredLogger
}

func (r *ReleaseLogger) Info(args ...interface{}) {
}

func (r *ReleaseLogger) Infof(template string, args ...interface{}) {
}

func NewLogger(mode config.Mode) Logger {

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	enc := zapcore.NewConsoleEncoder(cfg)

	enabler := zap.NewAtomicLevelAt(zap.DebugLevel)
	logger := zap.New(zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), enabler))
	zap.ReplaceGlobals(logger)

	log.SetFlags(0)
	log.SetOutput(Printer(logger.Debug))

	switch mode {
	case config.RELEASE:
		return &ReleaseLogger{SugaredLogger: logger.Sugar()}
	case config.DEBUG:
		fallthrough
	default:
		return logger.Sugar()
	}
}
