package handlers

import (
	"zatrasz75/gRPC_Interaction/auth/configs"
	"zatrasz75/gRPC_Interaction/auth/internal/repository"
	authPb "zatrasz75/gRPC_Interaction/auth/pkg/grpc/auth"
	"zatrasz75/gRPC_Interaction/auth/pkg/hash"
	"zatrasz75/gRPC_Interaction/auth/pkg/logger"
)

type Store struct {
	repo *repository.Store
	l    logger.LoggersInterface
	cfg  *configs.Config
	authPb.UnimplementedAuthServer
	h *hash.SGA1Hasher
}

func New(repo *repository.Store, l logger.LoggersInterface, cfg *configs.Config, haser *hash.SGA1Hasher) *Store {
	return &Store{
		repo: repo,
		l:    l,
		cfg:  cfg,
		h:    haser,
	}
}
