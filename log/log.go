package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// 로그 레벨
const (
	LevelDebug    = "debug"
	LevelInfo     = "info"
	LevelWarn     = "warn"
	LevelError    = "error"
	LevelCritical = "critical"
)

type Field = zapcore.Field

var Int64 = zap.Int64
var Int = zap.Int
var Uint32 = zap.Uint32
var String = zap.String
var Any = zap.Any
var Err = zap.Error
var Bool = zap.Bool

func getZapLevel(level string) zapcore.Level {
	switch level {
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// Logger 로그
type Logger struct {
	zap   *zap.Logger
	level zap.AtomicLevel
}

func NewLogger(name string, level string) *Logger {
	logger := &Logger{}

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename: name,
		MaxSize:  100, // 100mb 초과 시에 새로운 파일에 작성
		Compress: true,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writer,
		getZapLevel(level),
	)

	logger.zap = zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return logger
}

// Debug :
func (l *Logger) Debug(msg string, fields ...Field) {
	l.zap.Debug(msg, fields...)
}

// Info :
func (l *Logger) Info(msg string, fields ...Field) {
	l.zap.Info(msg, fields...)
}

// Warn :
func (l *Logger) Warn(msg string, fields ...Field) {
	l.zap.Warn(msg, fields...)
}

// Error :
func (l *Logger) Error(msg string, fields ...Field) {
	l.zap.Error(msg, fields...)
}

// Panic :
func (l *Logger) Panic(msg string, fields ...Field) {
	l.zap.Panic(msg, fields...)
}
