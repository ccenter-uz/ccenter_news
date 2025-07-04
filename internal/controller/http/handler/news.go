package handler

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mirjalilova/ccenter_news.git/internal/entity"
)

// CreateNews godoc
// @Summary Create a new News
// @Description Create a new News with the provided details
// @Tags News
// @Accept  json
// @Produce  json
// @Param Banner body entity.BannerCreate true "News Details"
// @Success 200 {object} entity.BannerCreate
// @Failure 400 {object}  string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /news/create [post]
func (h *Handler) CreateBanner(c *gin.Context) {
	reqBody := entity.BannerCreate{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	date, err := time.Parse("2006-01-02", reqBody.Date)
	if err != nil {
		c.JSON(400, gin.H{"Invalid date format (yyyy-mm-dd)": err})
		slog.Error("Invalid date format: ", "err", err, "date:", date)
		return
	}

	err = h.UseCase.BannerRepo.Create(context.Background(), &reqBody)
	if err != nil {
		c.JSON(500, gin.H{"Error creating Banner:": err})
		slog.Error("Error creating Banner: ", "err", err)
		return
	}

	slog.Info("New created successfully")
	c.JSON(200, gin.H{"Massage": "New created successfully"})
}

// GetByIdNews godoc
// @Summary Get News by ID
// @Description Get an News by their ID
// @Tags News
// @Accept  json
// @Produce  json
// @Param id query string true "News ID"
// @Success 200 {object} entity.BannerRes
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /news/get [get]
func (h *Handler) GetByIdBanner(c *gin.Context) {
	Banner_id := c.Query("id")

	res, err := h.UseCase.BannerRepo.GetById(context.Background(), &entity.ById{Id: Banner_id})
	if err != nil {
		c.JSON(500, gin.H{"Error getting Banner by ID: ": err})
		slog.Error("Error getting Banner by ID: ", "err", err)
		return
	}

	slog.Info("New retrieved successfully")
	c.JSON(200, res)
}

// UpdateNew godoc
// @Summary Update an New
// @Description Update an New's details
// @Tags News
// @Accept  json
// @Produce  json
// @Param id query string true "New ID"
// @Param Banner body entity.BannerCreate true "News Update Details"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /news/update [put]
func (h *Handler) UpdateBanner(c *gin.Context) {
	reqBody := entity.BannerCreate{}

	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body:": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	if reqBody.Date != "" && reqBody.Date != "string" {
		date, err := time.Parse("2006-01-02", reqBody.Date)
		if err != nil {
			c.JSON(400, gin.H{"Invalid date format (yyyy-mm-dd)": err})
			slog.Error("Invalid date format: ", "err", err, "date:", date)
			return
		}
	}

	err = h.UseCase.BannerRepo.Update(context.Background(), &entity.BannerUpdate{
		Id:       c.Query("id"),
		Text:     reqBody.Text,
		Title:    reqBody.Title,
		Label:    reqBody.Label,
		Date:     reqBody.Date,
		ImgUrl:   reqBody.ImgUrl,
		FileLink: reqBody.FileLink,
		HrefName: reqBody.HrefName,
		Type:     reqBody.Type,
		Order:    reqBody.Order,
	})
	if err != nil {
		c.JSON(500, gin.H{"Error updating New:": err})
		slog.Error("Error updating New: ", "err", err)
		return
	}

	slog.Info("New updated successfully")
	c.JSON(200, "New updated successfully")
}

// GetAllNews godoc
// @Summary Get all News
// @Description Get all News with optional filtering
// @Tags News
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} entity.BannerGetAllRes
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /news/list [get]
func (h *Handler) GetAllBanners(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")

	limitValue, offsetValue, err := parsePaginationParams(c, limit, offset)
	if err != nil {
		c.JSON(400, gin.H{"Error parsing pagination parameters:": err.Error()})
		slog.Error("Error parsing pagination parameters: ", "err", err)
		return
	}

	req := &entity.Filter{
		Limit:  limitValue,
		Offset: offsetValue,
	}

	res, err := h.UseCase.BannerRepo.GetAll(context.Background(), req)
	if err != nil {
		c.JSON(500, gin.H{"Error getting News:": err})
		slog.Error("Error getting News: ", "err", err)
		return
	}

	slog.Info("News retrieved successfully")
	c.JSON(200, res)
}

// DeleteNew godoc
// @Summary Delete an New
// @Description Delete an New by ID
// @Tags News
// @Accept  json
// @Produce  json
// @Param id query string true "New ID"
// @Success 200 {string} string "New deleted successfully"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /news/delete [delete]
func (h *Handler) DeleteBanner(c *gin.Context) {
	Banner_id := c.Query("id")

	err := h.UseCase.BannerRepo.Delete(context.Background(), &entity.ById{Id: Banner_id})
	if err != nil {
		c.JSON(500, gin.H{"Error deleting New by ID:": err})
		slog.Error("Error deleting New by ID: ", "err", err)
		return
	}

	slog.Info("New deleted successfully")
	c.JSON(200, "New deleted successfully")
}

// DeleteImage godoc
// @Summary Delete a Image
// @Description Delete a Image
// @Tags News
// @Accept  json
// @Produce  json
// @Param url query string true "Image url"
// @Success 200 {string} string "Image deleted successfully"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /news/image/delete [delete]
func (h *Handler) DeleteImage(c *gin.Context) {
	url := c.Query("url")

	err := h.UseCase.BannerRepo.DeleteImage(context.Background(), &entity.DeleteImage{ImgUrl: url})
	if err != nil {
		c.JSON(500, gin.H{"Error deleting Image:": err})
		slog.Error("Error deleting Image: ", "err", err)
		return
	}

	slog.Info("Image deleted successfully")
	c.JSON(200, "Image deleted successfully")
}

// ListImages godoc
// @Summary Get all Images
// @Description Get all Images
// @Tags News
// @Accept  json
// @Produce  json
// @Success 200 {object} []entity.Url
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /news/images/list [get]
func (h *Handler) ListImages(c *gin.Context) {

	res, err := h.UseCase.BannerRepo.GetImages(context.Background())
	if err != nil {
		c.JSON(500, gin.H{"Error getting images:": err})
		slog.Error("Error getting images: ", "err", err)
		return
	}

	urls := []entity.Url{}
	for _, url := range res.Images {
		if url.ImgUrl != "" {
			urls = append(urls, entity.Url{Url: url.ImgUrl})
		}
		if url.FileLink != "" {
			urls = append(urls, entity.Url{Url: url.FileLink})
		}
	}

	slog.Info("Images retrieved successfully")
	c.JSON(200, urls)
}

func parsePaginationParams(c *gin.Context, limit, offset string) (int, int, error) {
	limitValue := 10
	offsetValue := 0

	if limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err != nil {
			slog.Error("Invalid limit value", "err", err)
			c.JSON(400, gin.H{"error": "Invalid limit value"})
			return 0, 0, err
		}
		limitValue = parsedLimit
	} else {
		limitValue = 0
	}

	if offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err != nil {
			slog.Error("Invalid offset value", "err", err)
			c.JSON(400, gin.H{"error": "Invalid offset value"})
			return 0, 0, err
		}
		offsetValue = parsedOffset
	}

	return limitValue, offsetValue, nil
}
