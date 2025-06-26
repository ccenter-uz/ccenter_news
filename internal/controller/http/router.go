package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/mirjalilova/ccenter_news.git/config"
	_ "github.com/mirjalilova/ccenter_news.git/docs"
	"github.com/mirjalilova/ccenter_news.git/internal/controller/http/handler"
	"github.com/mirjalilova/ccenter_news.git/internal/usecase"
	"github.com/mirjalilova/ccenter_news.git/pkg/logger"
	"github.com/mirjalilova/ccenter_news.git/pkg/minio"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// NewRouter -.
// Swagger spec:
// @title       Ccenter News API
// @description This is a sample server Ccenter News server.
// @version     1.0
// @BasePath    /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(engine *gin.Engine, l *logger.Logger, config *config.Config, useCase *usecase.UseCase, minioClient *minio.MinIO) {
	// Options
	engine.Use(gin.Logger())
	//engine.Use(gin.Recovery())

	handlerV1 := handler.NewHandler(l, config, useCase, *minioClient)

	// Initialize Casbin enforcer

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Frontend domenini yozish
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Authentication"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	engine.Use(TimeoutMiddleware(5 * time.Second))
	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	// K8s probe
	engine.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	engine.Use(cors.Default())
	// Prometheus metrics
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	enforcer, err := casbin.NewEnforcer("./internal/controller/http/casbin/model.conf", "./internal/controller/http/casbin/policy.csv")
	if err != nil {
		slog.Error("Error while creating enforcer: ", "err", err)
	}

	if enforcer == nil {
		slog.Error("Enforcer is nil after initialization!")
	} else {
		slog.Info("Enforcer initialized successfully.")
	}

	// Routes

	// auth
	// engine.POST("/auth/login", handlerV1.Login)

	engine.POST("/img-upload", handlerV1.UploadFile)

	news := engine.Group("/news")
	{

		news.GET("/get", handlerV1.GetByIdBanner)
		news.GET("/list", handlerV1.GetAllBanners)
		news.POST("/create", handlerV1.CreateBanner)
		news.PUT("/update", handlerV1.UpdateBanner)
		news.DELETE("/delete", handlerV1.DeleteBanner)
		news.DELETE("/image/delete", handlerV1.DeleteImage)
	}
}
