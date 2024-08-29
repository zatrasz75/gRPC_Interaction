package gRPCserver

import (
	"net"
)

type Option func(*Server)

func OptionSet(host, port string) Option {
	return func(s *Server) {
		s.Addr = net.JoinHostPort(host, port)
	}
}
