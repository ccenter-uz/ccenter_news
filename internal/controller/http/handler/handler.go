package handler

import (
	"github.com/mirjalilova/ccenter_news.git/config"
	"github.com/mirjalilova/ccenter_news.git/internal/usecase"
	"github.com/mirjalilova/ccenter_news.git/pkg/logger"
)

type Handler struct {
	Logger  *logger.Logger
	Config  *config.Config
	UseCase *usecase.UseCase
	// MinIO   *minio.MinIO
}

func NewHandler(l *logger.Logger, c *config.Config, useCase *usecase.UseCase) *Handler {
	return &Handler{
		Logger:  l,
		Config:  c,
		UseCase: useCase,
		// MinIO:   &mn,
	}
}
