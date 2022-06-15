package log

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type correlationIDType int

const (
	requestIDKey correlationIDType = iota
)

var log *zap.Logger

//InitializeLogger Initialize logger
func InitializeLogger() {
	pe := zap.NewDevelopmentEncoderConfig()
	pe.ConsoleSeparator = " | "
	pe.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	consoleEncoder := zapcore.NewConsoleEncoder(pe)
	level := zap.InfoLevel

	core := zapcore.NewTee(zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))

	//Logs caller file name and skips 1 call stack
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// WithRqID returns a context which knows its request ID
func WithRqID(ctx context.Context, rqID string) context.Context {
	return context.WithValue(ctx, requestIDKey, rqID)
}

// Logger returns a zap logger with as much context as possible
func Logger(ctx context.Context) *zap.Logger {
	newLogger := log
	if ctx != nil {
		if ctxRqID, ok := ctx.Value(requestIDKey).(string); ok {
			newLogger = newLogger.With(zap.String("traceId", ctxRqID))
		}
	}
	return newLogger
}

//Info Logs message with Info level
func Info(ctx context.Context, msg string, value ...interface{}) {
	Logger(ctx).Sugar().Infof(fmt.Sprintf(msg, value...))
}

//Warn Logs message with Warn level
func Warn(ctx context.Context, msg string, value ...interface{}) {
	Logger(ctx).Sugar().Warnf(fmt.Sprintf(msg, value...))
}

//Error Logs message with Error level
func Error(ctx context.Context, msg string, value ...interface{}) {
	Logger(ctx).Sugar().Errorf(fmt.Sprintf(msg, value...))
}

//NewCtx Initialize new context with traceID
func NewCtx() context.Context {
	rqID, _ := uuid.NewRandom()
	return WithRqID(context.Background(), rqID.String())
}
