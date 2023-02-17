package main

import (
	"fmt"

	"go.uber.org/zap"
)

func initLogger(debug bool, logFile string) (*zap.SugaredLogger, error) {
	var zl *zap.Logger
	cfg := zap.Config{
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{logFile},
		ErrorOutputPaths: []string{logFile},
	}
	if debug {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	zl, err := cfg.Build()
	defer zl.Sync()

	if err != nil {
		return nil, fmt.Errorf("error when initializing logger: %w", err)
	}

	sugar := zl.Sugar()
	return sugar, nil
}
