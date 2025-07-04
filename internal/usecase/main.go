package usecase

import (
	"github.com/mirjalilova/ccenter_news.git/config"
	"github.com/mirjalilova/ccenter_news.git/internal/usecase/repo"
	"github.com/mirjalilova/ccenter_news.git/pkg/logger"
	"github.com/mirjalilova/ccenter_news.git/pkg/postgres"
)

type UseCase struct {
	BannerRepo BannerRepoI
}

func New(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		BannerRepo: repo.NewBannerRepo(pg, config, logger),
	}
}
