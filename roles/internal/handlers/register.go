package handlers

import (
	"zatrasz75/gRPC_Interaction/roles/configs"
	"zatrasz75/gRPC_Interaction/roles/internal/repository"
	rplesPb "zatrasz75/gRPC_Interaction/roles/pkg/grpc/roles"
	"zatrasz75/gRPC_Interaction/roles/pkg/logger"
)

type Store struct {
	repo *repository.Store
	l    logger.LoggersInterface
	cfg  *configs.Config
	rplesPb.UnimplementedRolesServer
}

func New(repo *repository.Store, l logger.LoggersInterface, cfg *configs.Config) *Store {
	return &Store{
		repo: repo,
		l:    l,
		cfg:  cfg,
	}
}
