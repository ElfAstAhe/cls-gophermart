package logger

import (
	_cfg "github.com/ElfAstAhe/cls-gophermart.git/internal/app/config"
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

func Initialize(level string, stage string) error {
	zapLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	config := getConfigByStage(stage)
	config.Level = zapLevel

	zl, err := config.Build()
	if err != nil {
		return err
	}

	Log = zl

	return nil
}

func getConfigByStage(stage string) *zap.Config {
	var config zap.Config
	switch stage {
	case _cfg.ProjectStageProduction:
		config = zap.NewProductionConfig()
	default:
		config = zap.NewDevelopmentConfig()
	}

	return &config
}
