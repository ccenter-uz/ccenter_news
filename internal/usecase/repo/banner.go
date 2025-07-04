package repo

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/mirjalilova/ccenter_news.git/config"
	"github.com/mirjalilova/ccenter_news.git/internal/entity"
	"github.com/mirjalilova/ccenter_news.git/pkg/logger"
	"github.com/mirjalilova/ccenter_news.git/pkg/postgres"
)

type BannerRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewBannerRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *BannerRepo {
	return &BannerRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *BannerRepo) Create(ctx context.Context, req *entity.BannerCreate) error {
	tx, err := r.pg.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if req.Order == 0 {
		var maxOrder int
		err := tx.QueryRow(ctx, `SELECT COALESCE(MAX("order"), 0) FROM banner`).Scan(&maxOrder)
		if err != nil {
			return err
		}
		req.Order = maxOrder + 1
	} else {
		_, err := tx.Exec(ctx, `
			UPDATE banner
			SET "order" = "order" + 1
			WHERE "order" >= $1
		`, req.Order)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO banner (
			text_uz,
			text_ru,
			text_en,
			title_uz,
			title_ru,
			title_en,
			date,
			label_uz,
			label_ru,
			label_en,
			img_url,
			file_link,
			href_name,
			type,
			"order"
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	_, err = tx.Exec(ctx, query,
		req.Text.Uz, req.Text.Ru, req.Text.En,
		req.Title.Uz, req.Title.Ru, req.Title.En,
		req.Date,
		req.Label.Uz, req.Label.Ru, req.Label.En,
		req.ImgUrl, req.FileLink, req.HrefName,
		req.Type, req.Order,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *BannerRepo) GetById(ctx context.Context, req *entity.ById) (*entity.BannerRes, error) {

	var res entity.BannerRes
	var createdAt time.Time

	query := `
	SELECT
		id,
		text_uz,
		text_ru,
		text_en,
		title_uz,
		title_ru,
		title_en,
		label_uz,
		label_ru,
		label_en,
		img_url,
		file_link,
		date,
		href_name,
  		type,
		"order",
		created_at
	FROM 
		banner
	WHERE 
		deleted_at = 0
	AND 
		id = $1
	`

	row := r.pg.Pool.QueryRow(ctx, query, req.Id)
	err := row.Scan(
		&res.Id,
		&res.Text.Uz,
		&res.Text.Ru,
		&res.Text.En,
		&res.Title.Uz,
		&res.Title.Ru,
		&res.Title.En,
		&res.Label.Uz,
		&res.Label.Ru,
		&res.Label.En,
		&res.ImgUrl,
		&res.FileLink,
		&res.Date,
		&res.HrefName,
		&res.Type,
		&res.Order,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}
	res.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

	return &res, nil
}

func (r *BannerRepo) GetAll(ctx context.Context, req *entity.Filter) (*entity.BannerGetAllRes, error) {

	resp := &entity.BannerGetAllRes{}

	query := `
	SELECT
		COUNT(id) OVER () AS total_count,
		id,
		text_uz,
		text_ru,
		text_en,
		title_uz,
		title_ru,
		title_en,
		label_uz,
		label_ru,
		label_en,
		img_url,
		file_link,
		date,
		href_name,
		type,
		"order",
		created_at
	FROM
		banner
	WHERE 
		deleted_at = 0 ORDER BY "order"
	`

	var args []interface{}

	if req.Limit == 0 {
		query += " OFFSET $1"
		args = append(args, req.Offset)
	} else {
		query += " LIMIT $1 OFFSET $2"
		args = append(args, req.Limit)
		args = append(args, req.Offset)
	}
	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("banner list is empty")
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		res := entity.BannerRes{}
		var count int
		var createdAt time.Time

		err := rows.Scan(
			&count,
			&res.Id,
			&res.Text.Uz,
			&res.Text.Ru,
			&res.Text.En,
			&res.Title.Uz,
			&res.Title.Ru,
			&res.Title.En,
			&res.Label.Uz,
			&res.Label.Ru,
			&res.Label.En,
			&res.ImgUrl,
			&res.FileLink,
			&res.Date,
			&res.HrefName,
			&res.Type,
			&res.Order,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		res.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

		resp.Banners = append(resp.Banners, res)
		resp.Count = count
	}

	return resp, nil
}

func (r *BannerRepo) Update(ctx context.Context, req *entity.BannerUpdate) error {
	query := `
	UPDATE
		banner
	SET`

	var conditions []string
	var args []interface{}

	if req.Text.En != "" && req.Text.En != "string" {
		conditions = append(conditions, " text_en = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Text.En)
	}
	if req.Text.Ru != "" && req.Text.Ru != "string" {
		conditions = append(conditions, " text_ru = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Text.Ru)
	}
	if req.Text.Uz != "" && req.Text.Uz != "string" {
		conditions = append(conditions, " text_uz = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Text.Uz)
	}
	if req.Title.En != "" && req.Title.En != "string" {
		conditions = append(conditions, " title_en = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Title.En)
	}
	if req.Title.Ru != "" && req.Title.Ru != "string" {
		conditions = append(conditions, " title_ru = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Title.Ru)
	}
	if req.Title.Uz != "" && req.Title.Uz != "string" {
		conditions = append(conditions, " title_uz = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Title.Uz)
	}
	if req.Label.En != "" && req.Label.En != "string" {
		conditions = append(conditions, " label_en = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Label.En)
	}
	if req.Label.Ru != "" && req.Label.Ru != "string" {
		conditions = append(conditions, " label_ru = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Label.Ru)
	}
	if req.Label.Uz != "" && req.Label.Uz != "string" {
		conditions = append(conditions, " label_uz = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Label.Uz)
	}
	if req.Date != "" && req.Date != "string" {
		conditions = append(conditions, " date = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Date)
	}
	if req.ImgUrl != "" && req.ImgUrl != "string" {
		conditions = append(conditions, " img_url = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.ImgUrl)
	}
	if req.FileLink != "" && req.FileLink != "string" {
		conditions = append(conditions, " file_link = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.FileLink)
	}
	if req.HrefName != "" && req.HrefName != "string" {
		conditions = append(conditions, " href_name = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.HrefName)
	}
	if req.Type != "" && req.Type != "string" {
		conditions = append(conditions, " type = $"+strconv.Itoa(len(args)+1))
		args = append(args, req.Type)
	}
	if req.Order != 0 {
		err := UpdateBannerOrder(ctx, r.pg, req.Id, req.Order)
		if err != nil {
			return err
		}
	}

	conditions = append(conditions, " updated_at = now()")
	query += strings.Join(conditions, ", ")
	query += " WHERE id = $" + strconv.Itoa(len(args)+1) + " AND deleted_at = 0"

	args = append(args, req.Id)

	_, err := r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *BannerRepo) Delete(ctx context.Context, req *entity.ById) error {
	tx, err := r.pg.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var order int
	err = tx.QueryRow(ctx, `SELECT "order" FROM banner WHERE id = $1 AND deleted_at = 0`, req.Id).Scan(&order)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE banner SET "order" = -(EXTRACT(EPOCH FROM NOW())) WHERE id = $1`, req.Id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE banner
		SET "order" = "order" - 1
		WHERE "order" > $1
	`, order)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE banner
		SET deleted_at = EXTRACT(EPOCH FROM NOW())
		WHERE id = $1
	`, req.Id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}



func (r *BannerRepo) DeleteImage(ctx context.Context, req *entity.DeleteImage) error {

	query := `
	UPDATE 
		banner
	SET 
		img_url = ''
	WHERE 
		img_url = $1
	AND 
		deleted_at = 0
	`

	_, err := r.pg.Pool.Exec(ctx, query, req.ImgUrl)

	if err != nil {
		return err
	}

	return nil
}

func (r *BannerRepo) GetImages(ctx context.Context) (*entity.ListImages, error) {

	resp := &entity.ListImages{}

	query := `
	SELECT
		COUNT(id) OVER () AS total_count,
		img_url,
		file_link
	FROM
		banner
	WHERE 
		deleted_at = 0 AND img_url <> ''
	`

	rows, err := r.pg.Pool.Query(ctx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("imgage not found")
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		res := entity.Image{}
		var count int

		err := rows.Scan(
			&count,
			&res.ImgUrl,
			&res.FileLink,
		)
		if err != nil {
			return nil, err
		}

		resp.Images = append(resp.Images, res)
		resp.Count = count
	}

	return resp, nil
}

func UpdateBannerOrder(ctx context.Context, db *postgres.Postgres, id string, newOrder int) error {
	var oldOrder int

	err := db.Pool.QueryRow(ctx, `SELECT "order" FROM banner WHERE id = $1`, id).Scan(&oldOrder)
	if err != nil {
		return err
	}

	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	defer tx.Rollback(ctx)

	if newOrder < oldOrder {
		_, err = tx.Exec(ctx, `
			UPDATE banner
			SET "order" = "order" + 1
			WHERE "order" >= $1 AND "order" < $2
		`, newOrder, oldOrder)
	} else if newOrder > oldOrder {
		_, err = tx.Exec(ctx, `
			UPDATE banner
			SET "order" = "order" - 1
			WHERE "order" <= $1 AND "order" > $2
		`, newOrder, oldOrder)
	}
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE banner
		SET "order" = $1
		WHERE id = $2
	`, newOrder, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
