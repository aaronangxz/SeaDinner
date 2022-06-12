package Log

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type correlationIdType int

const (
	requestIdKey correlationIdType = iota
	sessionIdKey
)

var Log *zap.Logger

func InitializeLogger() {
	pe := zap.NewDevelopmentEncoderConfig()
	pe.ConsoleSeparator = " | "
	pe.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	consoleEncoder := zapcore.NewConsoleEncoder(pe)
	level := zap.InfoLevel

	core := zapcore.NewTee(zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))

	//Logs caller file name and skips 1 call stack
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// WithRqId returns a context which knows its request ID
func WithRqId(ctx context.Context, rqId string) context.Context {
	return context.WithValue(ctx, requestIdKey, rqId)
}

// Logger returns a zap logger with as much context as possible
func Logger(ctx context.Context) *zap.Logger {
	newLogger := Log
	if ctx != nil {
		if ctxRqId, ok := ctx.Value(requestIdKey).(string); ok {
			newLogger = newLogger.With(zap.String("traceId", ctxRqId))
		}
	}
	return newLogger
}

func Info(ctx context.Context, msg string, value ...interface{}) {
	Logger(ctx).Sugar().Infof(fmt.Sprintf(msg, value...))
}

func Warn(ctx context.Context, msg string, value ...interface{}) {
	Logger(ctx).Sugar().Warnf(fmt.Sprintf(msg, value...))
}

func Error(ctx context.Context, msg string, value ...interface{}) {
	Logger(ctx).Sugar().Errorf(fmt.Sprintf(msg, value...))
}

func NewCtx() context.Context {
	rqId, _ := uuid.NewRandom()
	return WithRqId(context.Background(), rqId.String())
}
