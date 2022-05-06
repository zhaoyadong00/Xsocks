package logs

import (
	"Xsocks-core/util/config"
	"go.uber.org/zap"
)

func init() {
	InitNewZap()
}

var Logger *zap.Logger

func InitNewZap()  {
	if config.Config["level"] == "debug" {
		Logger , _ = zap.NewDevelopment()
	}else{
		Logger , _ = zap.NewProduction()
	}
	defer Logger.Sync()
}
