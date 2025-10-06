package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapLogger struct {
	logger *zap.Logger
}

func NewStartupZapLogger() AppLogger {
	zapLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	core := newConsoleZapCore(zapLevel)

	return &zapLogger{
		logger: zap.New(*core),
	}
}

func NewZapLogger(level string, filePath string) (AppLogger, error) {
	zapLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	consoleCore := newConsoleZapCore(zapLevel)
	core := zapcore.NewTee(*consoleCore)
	if strings.TrimSpace(filePath) != "" {
		fileCore := newFileZapCore(zapLevel, filePath)

		core = zapcore.NewTee(*consoleCore, *fileCore)
	}

	res := &zapLogger{}

	res.logger = zap.New(core)

	return res, nil
}

func newConsoleZapCore(level zap.AtomicLevel) *zapcore.Core {
	stdOut := zapcore.AddSync(os.Stdout)

	encConf := zap.NewDevelopmentEncoderConfig()
	encConf.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(encConf)

	res := zapcore.NewCore(consoleEncoder, stdOut, level)

	return &res
}

func newFileZapCore(level zap.AtomicLevel, filePath string) *zapcore.Core {
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
	})

	encConf := zap.NewProductionEncoderConfig()
	encConf.TimeKey = "timestamp"
	encConf.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(encConf)

	res := zapcore.NewCore(fileEncoder, file, level)

	return &res
}

// Closer

func (zl *zapLogger) Close() error {
	return zl.logger.Sync()
}

// AppLogger

func (zl *zapLogger) Error(args ...interface{}) {
	zl.logger.Sugar().Error(args...)
}

func (zl *zapLogger) Errorf(format string, args ...interface{}) {
	zl.logger.Sugar().Errorf(format, args...)
}

func (zl *zapLogger) Warn(args ...interface{}) {
	zl.logger.Sugar().Warn(args...)
}

func (zl *zapLogger) Warnf(format string, args ...interface{}) {
	zl.logger.Sugar().Warnf(format, args...)
}

func (zl *zapLogger) Info(args ...interface{}) {
	zl.logger.Sugar().Info(args...)
}

func (zl *zapLogger) Infof(format string, args ...interface{}) {
	zl.logger.Sugar().Infof(format, args...)
}

func (zl *zapLogger) Debug(args ...interface{}) {
	zl.logger.Sugar().Debug(args...)
}

func (zl *zapLogger) Debugf(format string, args ...interface{}) {
	zl.logger.Sugar().Debugf(format, args...)
}

func (zl *zapLogger) GetLogger(logicEntry string) AppLogger {
	return &zapLogger{
		logger: zl.logger.With(zap.String("childEntry", logicEntry)),
	}
}
