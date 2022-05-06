package socks5

import "net"

type Server struct {
	Config *ServerConfig
}

func NewServer(config *ServerConfig) *Server {
	return &Server{
		Config: config,
	}
}

func (s *Server)Process() {
	//net.Listen()

	net.ListenTCP("tcp" , s.Config.Address)
}
