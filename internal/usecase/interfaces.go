// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/mirjalilova/ccenter_news.git/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// AuthRepo -.
	AuthRepoI interface {
		Login(ctx context.Context, req *entity.LoginReq) (*entity.UserInfo, error)
	}

	// BannerRepo -.
	BannerRepoI interface {
		Create(ctx context.Context, req *entity.BannerCreate) error
		GetById(ctx context.Context, req *entity.ById) (*entity.BannerRes, error)
		GetAll(ctx context.Context, req *entity.Filter) (*entity.BannerGetAllRes, error)
		Update(ctx context.Context, req *entity.BannerUpdate) error
		Delete(ctx context.Context, req *entity.ById) error
		DeleteImage(ctx context.Context, req *entity.DeleteImage) error
		GetImages(ctx context.Context) (*entity.ListImages, error)
	}
)
