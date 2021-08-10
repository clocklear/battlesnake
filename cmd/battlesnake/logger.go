package main

import (
	"os"

	"github.com/go-kit/kit/log"
)

type logger struct {
	base log.Logger
}

func (l logger) LogWithLevel(level string, message string, keyvals ...interface{}) {
	params := []interface{}{
		"level", level,
		"msg", message,
	}
	params = append(params, keyvals...)
	_ = l.base.Log(params...)
}

func (l logger) Info(message string, keyvals ...interface{}) {
	l.LogWithLevel("info", message, keyvals...)
}

func (l logger) Error(message string, keyvals ...interface{}) {
	l.LogWithLevel("error", message, keyvals...)
}

func (l logger) Fatal(message string, keyvals ...interface{}) {
	l.LogWithLevel("fatal", message, keyvals...)
	os.Exit(1)
}
