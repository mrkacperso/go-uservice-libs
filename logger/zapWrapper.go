package logger

import (
	"fmt"
	"gitlab.com/kodeit/turbotrader-libs/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func SetupLogger(timeFormat string, configMode int) (*zap.Logger, error) {
	var zapCfg zap.Config
	switch configMode {
	case config.Dev:
		zapCfg = zap.NewDevelopmentConfig()
	case config.Prod:
		zapCfg = zap.NewProductionConfig()
	default:
		return nil, fmt.Errorf("invalid config mode %d", configMode)
	}

	//Setup custom time logger with time format
	zapCfg.EncoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format(timeFormat))
	}
	l, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("error initializing logger: %v", err)
	}
	return l, nil
}
