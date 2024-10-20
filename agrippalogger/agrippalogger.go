package agrippalogger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Core logger for the agrippa package
func GetLogger() *zap.SugaredLogger {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleSink := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(consoleEncoder, consoleSink, zap.
		InfoLevel)
	logger := zap.New(core)
	sugar := logger.Sugar()
	sugar.Debugw("logger init")
	return sugar
}
