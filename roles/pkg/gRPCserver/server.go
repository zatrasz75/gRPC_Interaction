package gRPCserver

import (
	"google.golang.org/grpc"
	"net"
	"time"
)

type Server struct {
	GrpcServer      *grpc.Server
	notify          chan error
	shutdownTimeout time.Duration
	Addr            string
}

func New(unaryInterceptors []grpc.UnaryServerInterceptor, options ...Option) *Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
	)

	s := &Server{
		GrpcServer: grpcServer,
		notify:     make(chan error, 1),
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	go func() {
		if err = s.GrpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			s.notify <- err
		}
	}()

	return nil
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() {
	s.GrpcServer.GracefulStop()
}
