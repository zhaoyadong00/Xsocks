package socks

import "net"

type Server struct {
	UserName          string
	Password          string
	Method            byte
	TCPAddr           *net.TCPAddr
	TCPListen         *net.TCPListener
	TCPTimeout        int
}

func NewServer() (*Server, error) {
	return &Server{}, nil
}
