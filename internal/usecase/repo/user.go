package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mirjalilova/ccenter_news.git/config"
	"github.com/mirjalilova/ccenter_news.git/internal/entity"
	"github.com/mirjalilova/ccenter_news.git/pkg/logger"
	"github.com/mirjalilova/ccenter_news.git/pkg/postgres"
)

type AuthRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewAuthRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *AuthRepo {
	return &AuthRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *AuthRepo) Login(ctx context.Context, req *entity.LoginReq) (*entity.UserInfo, error) {
	res := &entity.UserInfo{}
	
	var password string
	var createdAt time.Time
	query := `SELECT id, login, role, password_hash, service_name, username, first_number, image_url, created_at FROM users WHERE login = $1 AND password_hash = $2 AND deleted_at = 0`
	err := r.pg.Pool.QueryRow(ctx, query, req.Login, req.Password).Scan(
		&res.AgentID,
		&res.Login,
		&res.Role,
		&password,
		&res.ServiceName,
		&res.Name,
		&res.FirstNumber,
		&res.Image,
		&createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)); err != nil {
	// 	return nil, errors.New("invalid login or password")
	// }

	res.CreateDate = createdAt.Format("2006-01-02 15:04:05")

	return res, nil
}
