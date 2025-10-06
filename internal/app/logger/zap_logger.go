package logger

import (
	_cfg "github.com/ElfAstAhe/cls-gophermart.git/internal/app/config"
	"go.uber.org/zap"
)

func NewZapLogger(level string, stage string) (*zap.Logger, error) {
	zapLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	config := getZapConfigByStage(stage)
	config.Level = zapLevel

	zl, err := config.Build()
	if err != nil {
		return nil, err
	}

	return zl, nil
}

func getZapConfigByStage(stage string) *zap.Config {
	var config zap.Config
	switch stage {
	case _cfg.ProjectStageProduction:
		config = zap.NewProductionConfig()
	default:
		config = zap.NewDevelopmentConfig()
	}

	return &config
}
