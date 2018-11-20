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

var globalLogger *zap.Logger

// Field : 필드
type Field = zapcore.Field

// 필드를 반환하는 함수.
// ex) StringField(key, value)
var (
	Int64Field  = zap.Int64
	IntField    = zap.Int
	Uint32Field = zap.Uint32
	StringField = zap.String
	AnyField    = zap.Any
	ErrorField  = zap.Error
	BoolField   = zap.Bool
)

// 로그 레벨
const (
	DebugLevel  = "debug"
	InfoLevel   = "info"
	WarnLevel   = "warn"
	ErrorLevel  = "error"
	DPanicLevel = "dpanic"
	PanicLevel  = "panic"
	FatalLevel  = "fatal"
)

func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// InitGolbalLogger : 로그 설정을 초기화.
//
// 로그 레벨은 설정된 로그 레벨 이상 로그파일에 작성
func InitGolbalLogger(name string, level string) {
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename: name,
		MaxSize:  100, // 100mb 초과 시에 새로운 파일에 작성
		Compress: true,
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		writer,
		getZapLevel(level),
	)

	globalLogger = zap.New(core)
}

// Debug 로그 생성
func Debug(msg string, fields ...Field) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		globalLogger.Debug(fmt.Sprintf("%s:%d:%s", file, line, msg), fields...)
	} else {
		globalLogger.Debug(msg, fields...)
	}
}

// Info 로그 생성
func Info(msg string, fields ...Field) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		globalLogger.Info(fmt.Sprintf("%s:%d:%s", file, line, msg), fields...)
	} else {
		globalLogger.Info(msg, fields...)
	}
}

// Warn 로그 생성
func Warn(msg string, fields ...Field) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		globalLogger.Warn(fmt.Sprintf("%s:%d:%s", file, line, msg), fields...)
	} else {
		globalLogger.Warn(msg, fields...)
	}
}

// Error 로그 생성
func Error(msg string, fields ...Field) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		globalLogger.Error(fmt.Sprintf("%s:%d:%s", file, line, msg), fields...)
	} else {
		globalLogger.Error(msg, fields...)
	}
}

// Fatal 로그 생성
func Fatal(msg string, fields ...Field) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		globalLogger.Fatal(fmt.Sprintf("%s:%d:%s", file, line, msg), fields...)
	} else {
		globalLogger.Fatal(msg, fields...)
	}
}

// Panic 로그 생성
func Panic(msg string, fields ...Field) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		globalLogger.Panic(fmt.Sprintf("%s:%d:%s", file, line, msg), fields...)
	} else {
		globalLogger.Panic(msg, fields...)
	}
}
