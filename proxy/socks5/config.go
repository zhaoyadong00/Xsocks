package socks5

import (
	"Xsocks-core/util/logs"
	"github.com/goinggo/mapstructure"
	"go.uber.org/zap"
	"net"
)

type AuthType int32

const (
	// NO_AUTH is for anounymous authentication.
	AuthType_NO_AUTH AuthType = 0
	// PASSWORD is for username/password authentication.
	AuthType_PASSWORD AuthType = 1
)

type Version int32

const (
	Version_SOCKS5  Version = 0
	Version_SOCKS4  Version = 1
	Version_SOCKS4A Version = 2
)

type ServerConfig struct {
	AuthType AuthType
	Accounts map[string]string
	Address  *net.TCPAddr
	Version  Version
	Timeout  int32
}

type ClientConfig struct {
	Address net.TCPAddr
	Version Version
}

func GetServerConfig(config map[string]string) (*ServerConfig , error) {
	sc := &ServerConfig{}
	err := mapstructure.Decode(config , sc)
	if err != nil {
		logs.Logger.Info("mapstructure.Decode" , zap.Error(err))
	}
	return sc , err
}