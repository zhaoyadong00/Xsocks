package Xsocks

import (
	"Xsocks-core/proxy/socks5"
	"Xsocks-core/util/config"
	"Xsocks-core/util/logs"
	"go.uber.org/zap"
)

func Start()  {
	logs.Logger.Info("start Xsocks...............")
	logs.Logger.Info("config.config" , zap.Any("config",config.Config))
	sc , err := socks5.GetServerConfig(config.Config["inbound"].(map[string]string))
	if err != nil {
		logs.Logger.Panic("server config err:" , zap.Error(err))
	}
	srv := socks5.NewServer(sc)
	srv.Process()
}