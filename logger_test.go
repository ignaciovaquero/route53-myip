package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitSugaredLogger(t *testing.T) {
	cfg := zap.Config{
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
	}
	nonVerbose, _ := cfg.Build()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	verbose, _ := cfg.Build()

	testCases := []struct {
		name     string
		debug    bool
		expected *zap.SugaredLogger
	}{
		{
			name:     "Non verbose",
			debug:    false,
			expected: nonVerbose.Sugar(),
		},
		{
			name:     "Verbose",
			debug:    false,
			expected: verbose.Sugar(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			actual, _ := initLogger(tc.debug, "stderr")
			assert.ObjectsAreEqual(tc.expected, actual)
		})
	}
}
