package utils

import (
	"treehole_next/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLog() (*zap.Logger, error) {
	var atomicLevel zapcore.Level
	if config.Config.Mode != "production" {
		atomicLevel = zapcore.DebugLevel
	} else {
		atomicLevel = zapcore.InfoLevel
	}
	logConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(atomicLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	return logConfig.Build()
}

type Role string

const (
	RoleOwner    = "owner"
	RoleAdmin    = "admin"
	RoleOperator = "operator"
)

func MyLog(model string, action string, objectID, userID int, role Role, msg ...string) {
	message := ""
	for _, v := range msg {
		message += v
	}
	Logger.Info(
		message,
		zap.String("model", model),
		zap.Int("user_id", userID),
		zap.Int("object_id", objectID),
		zap.String("action", action),
		zap.String("role", string(role)),
	)
}
