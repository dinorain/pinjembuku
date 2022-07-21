package service

import (
	"github.com/dinorain/pinjembuku/config"
	"github.com/dinorain/pinjembuku/internal/session"
	"github.com/dinorain/pinjembuku/internal/user"
	"github.com/dinorain/pinjembuku/pkg/logger"
)

type usersServiceGRPC struct {
	logger logger.Logger
	cfg    *config.Config
	userUC user.UserUseCase
	sessUC session.SessUseCase
}

// Auth service constructor
func NewAuthServerGRPC(logger logger.Logger, cfg *config.Config, userUC user.UserUseCase, sessUC session.SessUseCase) *usersServiceGRPC {
	return &usersServiceGRPC{logger: logger, cfg: cfg, userUC: userUC, sessUC: sessUC}
}
