package migrations

import "go.uber.org/zap"

// gooseLogger is implementation of goose.Logger interface
type gooseLogger struct {
	log *zap.SugaredLogger
}

func newGooseLogger(log *zap.SugaredLogger) *gooseLogger {
	return &gooseLogger{
		log: log,
	}
}

// Logger

func (g *gooseLogger) Fatalf(format string, v ...interface{}) {
	g.log.Fatalf(format, v...)
}

func (g *gooseLogger) Printf(format string, v ...interface{}) {
	g.log.Infof(format, v...)
}

// ==========
