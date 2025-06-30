package api

import (
	"fmt"
	"golang-workshop/src/api/middleware"
	"golang-workshop/src/api/router"
	"golang-workshop/src/api/validation"
	"golang-workshop/src/config"
	"golang-workshop/src/docs"
	"golang-workshop/src/pkg/logging"

	_ "golang-workshop/src/docs"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitServer(cfg *config.Config, logger *logging.ZapLogger) {
	gin.SetMode(cfg.Server.RunMode)
	r := gin.New()
	RegisterValidators(logger)

	r.Use(middleware.Cors(cfg))
	r.Use(gin.Logger(), gin.CustomRecovery(middleware.ErrorHandler))

	RegisterRoutes(r, cfg)
	RegisterSwagger(r, cfg)
	logger.Info(logging.General, logging.Startup, "Started", nil)
	err := r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		logger.Fatal(logging.General, logging.Startup, err.Error(), nil)
	}
}

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	api := r.Group("/api")

	v1 := api.Group("/v1")
	{
		health := v1.Group("/health")
		router.Health(health)
	}

}

func RegisterValidators(logger *logging.ZapLogger) {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		err := val.RegisterValidation("mobile", validation.IranianMobileNumberValidator, true)
		if err != nil {
			logger.Error(logging.Validation, logging.Startup, err.Error(), nil)
		}
	}
}

func RegisterSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = "golang web api"
	docs.SwaggerInfo.Description = "golang web api"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Server.Port)
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
