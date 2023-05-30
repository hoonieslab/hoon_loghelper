package loghelp

import (
	"bytes"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"time"
)

var consoleLogger *zap.Logger
var fileLogger *zap.Logger

/**
 * Escape String JSON Encoder : START
 */

type EscapeSeqJSONEncoder struct {
	zapcore.Encoder
}

func (enc *EscapeSeqJSONEncoder) Clone() zapcore.Encoder {
	return enc
}

func (enc *EscapeSeqJSONEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// Interface Embed 수행
	b, err := enc.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}
	newb := buffer.NewPool().Get()

	// 이스케이핑 할 문자 재정열
	replaceStr := bytes.Replace(b.Bytes(), []byte("\\n"), []byte("\n"), -1)
	replaceStr = bytes.Replace(replaceStr, []byte("\\r"), []byte("\r"), -1)
	replaceStr = bytes.Replace(replaceStr, []byte("\\t"), []byte("\t"), -1)

	_, _ = newb.Write(replaceStr)

	return newb, nil
}

/**
 * Escape String JSON Encoder : END
 */

// Init intialize hoon log
func Init() {
	logFile := *logFilePath
	rotator, err := rotatelogs.New(
		logFile,
		rotatelogs.WithMaxAge(3*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour))
	if err != nil {
		panic(err)
	}

	//로거 정의
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		TimeKey:        "timestamp",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	logCore := zapcore.NewCore(
		&EscapeSeqJSONEncoder{Encoder: zapcore.NewJSONEncoder(encoderCfg)},
		zapcore.AddSync(rotator),
		zap.InfoLevel)

	fileLogger = zap.New(logCore)

	//로거 정의
	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""
	config.EncoderConfig = encoderCfg
	config.Encoding = "console"

	consoleLogger, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

// Info logging info Lvl both console and file
func Info(message string, fields ...zap.Field) {
	consoleLogger.Info(message, fields...)
	fileLogger.Info(message, fields...)
}

// Debug logging debug Lvl both console and file
func Debug(message string, fields ...zap.Field) {
	consoleLogger.Debug(message, fields...)
	fileLogger.Debug(message, fields...)
}

// Warn logging warn Lvl both console and file
func Warn(message string, fields ...zap.Field) {
	consoleLogger.Warn(message, fields...)
	fileLogger.Warn(message, fields...)
}

// Panic logging panic Lvl both console and file
func Panic(message string, fields ...zap.Field) {
	consoleLogger.Panic(message, fields...)
	fileLogger.Panic(message, fields...)
}

// Error logging error Lvl both console and file
func Error(message string, fields ...zap.Field) {
	consoleLogger.Error(message, fields...)
	fileLogger.Error(message, fields...)
}

// Fatal logging fatal Lvl both console and file
func Fatal(message string, fields ...zap.Field) {
	consoleLogger.Fatal(message, fields...)
	fileLogger.Fatal(message, fields...)
}
