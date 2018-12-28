// Package util 은 file, log 관련 기능들을 제공한다.
//
// log는 내부적으로 zap, lumberjack.v2를 사용
package util

import (
	"fmt"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Field : 필드
type Field = zapcore.Field

// 필드를 반환하는 함수.
// ex) StringField(key, value)
var (
	Int64  = zap.Int64
	Int    = zap.Int
	Uint32 = zap.Uint32
	String = zap.String
	Any    = zap.Any
	Error  = zap.Error
	Bool   = zap.Bool
)

// 로그 레벨
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

// Logger :
type Logger struct {
	zap *zap.Logger
}

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

// NewLogger :
func NewLogger(name string, level string) *Logger {
	logger := &Logger{}

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename: name,
		MaxSize:  100, // 100mb 초과 시에 새로운 파일에 작성
		Compress: true,
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		writer,
		getZapLevel(level),
	)

	logger.zap = zap.New(core)

	return logger
}

// Debug 로그 생성
func (l *Logger) Debug(msg string, fields ...Field) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		l.zap.Debug(fmt.Sprintf("%s:%d:%s", file, line, msg), fields...)
	} else {
		l.zap.Debug(msg, fields...)
	}
}

// Info 로그 생성
func (l *Logger) Info(msg string, fields ...Field) {
	l.zap.Info(msg, fields...)
}

// Warn 로그 생성
func (l *Logger) Warn(msg string, fields ...Field) {
	l.zap.Warn(msg, fields...)
}

// Error 로그 생성
func (l *Logger) Error(msg string, fields ...Field) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		l.zap.Error(fmt.Sprintf("%s:%d:%s", file, line, msg), fields...)
	} else {
		l.zap.Error(msg, fields...)
	}
}

// Panic 로그 생성
func (l *Logger) Panic(msg string, fields ...Field) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		l.zap.Panic(fmt.Sprintf("%s:%d:%s", file, line, msg), fields...)
	} else {
		l.zap.Panic(msg, fields...)
	}
}
