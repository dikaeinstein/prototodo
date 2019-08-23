package logger

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger creates a new zap Logger.
func NewZapLogger(lvl int, timeFormat string, appEnv string) *zap.Logger {
	// First, define our level-handling logic.
	globalLevel := zapcore.Level(lvl)
	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	// It is usefull for Kubernetes deployment.
	// Kubernetes interprets os.Stdout log items as INFO and os.Stderr log items
	// as ERROR by default.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= globalLevel && lvl < zapcore.ErrorLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Optimize the console output for for human operators.
	eCfg := zap.NewProductionEncoderConfig()
	// customTimeEncoder encode Time to our custom format
	// This example how we can customize zap default functionality
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(timeFormat))
	}
	eCfg.EncodeTime = customTimeEncoder

	var consoleEncoder zapcore.Encoder
	// set consoleEncoder based on our APP_ENV
	if !strings.Contains(appEnv, "production") || !strings.Contains(appEnv, "prod") {
		consoleEncoder = zapcore.NewConsoleEncoder(eCfg)
	} else {
		consoleEncoder = zapcore.NewJSONEncoder(eCfg)
	}
	zapcore.NewConsoleEncoder(eCfg)

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	return zap.New(core)
}
