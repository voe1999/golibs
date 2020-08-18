package logger

import (
	"errors"
	"github.com/voe1999/golibs/cmd/logger/extensions"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type LogWriter interface {
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
}

func (l Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

func (l Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.logger.Panicw(msg, keysAndValues...)
}

func (l Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

func (l Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l Logger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l Logger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l Logger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l Logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l Logger) Panicf(template string, args ...interface{}) {
	l.logger.Panicf(template, args...)
}

func (l Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l Logger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

type Logger struct {
	logger *zap.SugaredLogger
}

// NewLogger 创建Logger。
// Option是日志格式、位置和等级的组合。
// skipDepth是Logger调用的层级。如果直接调用Logger就设为0，外层每封装一层加1。
func NewLogger(options []Option, skipDepth int) (*Logger, error) {
	var cores []zapcore.Core
	for _, op := range options {
		var encoder zapcore.Encoder
		switch op.OutputFormat {
		case FormatPlainText:
			encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		case FormatJSON:
			encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		default:
			return nil, errors.New("logger format不合法")
		}

		var writeTo zapcore.WriteSyncer
		switch op.WriteTo.Type {
		case WriteToStdout:
			writeTo = zapcore.AddSync(os.Stdout)
		case WriteToStderr:
			writeTo = zapcore.AddSync(os.Stderr)
		case WriteToFile:
			if op.WriteTo.Path == "" {
				return nil, errors.New("logger writeTo不合法")
			}
			lj := &lumberjack.Logger{
				Filename:   op.WriteTo.Path,
				MaxSize:    10, //MB
				MaxBackups: 0,
				MaxAge:     0, //days
			}
			writeTo = zapcore.AddSync(lj)
		case WriteToMessageQueue:
			// TODO 引入mq producer，实现Writer接口
			panic("not implemented")
		case WriteToQYWeiXinBot:
			writeTo = zapcore.AddSync(extensions.NewQyWeiXinWriter(op.WriteTo.Path))
		default:
			return nil, errors.New("logger writeTo不合法")
		}

		var priority zap.LevelEnablerFunc
		if op.MinLevel < _minLevel || op.MaxLevel > _maxLevel || op.MinLevel > op.MaxLevel {
			return nil, errors.New("logger level不合法: %+v")
		}
		priority = func(option Option) zap.LevelEnablerFunc {
			return func(lev zapcore.Level) bool {
				return lev <= zapcore.Level(option.MaxLevel) && lev >= zapcore.Level(option.MinLevel)
			}
		}(op)
		core := zapcore.NewCore(encoder, writeTo, priority)
		cores = append(cores, core)
	}

	core := zapcore.NewTee(cores...)
	if skipDepth < 0 {
		return nil, errors.New("logger skipDepth不合法")
	}

	lgr := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1+skipDepth),
		zap.AddStacktrace(zap.ErrorLevel))
	logger := lgr.Sugar()
	return &Logger{logger: logger}, nil
}
