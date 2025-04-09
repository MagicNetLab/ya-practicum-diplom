package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log = zap.NewNop()

// Init инициализация логера
func Init() error {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		zap.InfoLevel,
	)

	log = zap.New(core)

	return nil
}

// Info логирование сообщения уровня info
func Info(msg string, args ...LogArg) {
	log.Info(msg, prepareZapArgs(args)...)
}

// Error логирование сообщения уровня error
func Error(msg string, args ...LogArg) {
	log.Error(msg, prepareZapArgs(args)...)
}

// Fatal логирование сообщения уровня fatal
func Fatal(msg string, args ...LogArg) {
	log.Fatal(msg, prepareZapArgs(args)...)
}

// Debug логирование сообщения уровня debug
func Debug(msg string, args ...LogArg) {
	log.Debug(msg, prepareZapArgs(args)...)
}

// Warn логирование сообщения уровня warn
func Warn(msg string, args ...LogArg) {
	log.Warn(msg, prepareZapArgs(args)...)
}

// Close закрытие логера
func Close() error {
	return log.Sync()
}

// prepareZapArgs подготовка аргументов для логирования
func prepareZapArgs(args []LogArg) []zapcore.Field {
	var zapArgs []zapcore.Field
	for _, arg := range args {
		switch arg.argType {
		case "string":
			zapArgs = append(zapArgs, zap.String(arg.name, arg.value.(string)))
		case "int":
			zapArgs = append(zapArgs, zap.Int(arg.name, arg.value.(int)))
		case "duration":
			zapArgs = append(zapArgs, zap.Duration(arg.name, arg.value.(time.Duration)))
		}
	}

	return zapArgs
}
