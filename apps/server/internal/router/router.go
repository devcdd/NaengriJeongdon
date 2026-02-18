package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"naengrijeongdon/apps/server/internal/config"
	"naengrijeongdon/apps/server/internal/handlers"

	_ "naengrijeongdon/apps/server/docs"
)

func New(cfg config.Config, pool *pgxpool.Pool) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	hh := handlers.NewHealthHandler(pool)
	uh := handlers.NewUserHandler(cfg, pool)
	ssh := handlers.NewStorageSpaceHandler()
	fh := handlers.NewFoodHandler()
	ph := handlers.NewPantryItemHandler()
	mh := handlers.NewMetadataHandler()
	ah := handlers.NewAuthHandler(cfg, pool)

	api := r.Group("/api/v1")
	{
		health := api.Group("/health")
		{
			health.GET("/live", hh.Live)
			health.GET("/ready", hh.Ready)
		}

		users := api.Group("/users")
		{
			users.GET("/me", uh.GetMe)
			users.PATCH("/me", uh.UpdateMe)
		}

		spaces := api.Group("/storage-spaces")
		{
			spaces.GET("", ssh.List)
			spaces.POST("", ssh.Create)
			spaces.GET("/:spaceId", ssh.Get)
			spaces.PATCH("/:spaceId", ssh.Update)
			spaces.DELETE("/:spaceId", ssh.Delete)
		}

		foods := api.Group("/foods")
		{
			foods.GET("/autocomplete", fh.Autocomplete)
			foods.GET("/:foodId", fh.Get)
			foods.GET("/:foodId/predict-expiry", fh.PredictExpiry)
		}

		pantry := api.Group("/pantry-items")
		{
			pantry.GET("", ph.List)
			pantry.POST("", ph.Create)
			pantry.GET("/:itemId", ph.Get)
			pantry.PATCH("/:itemId", ph.Update)
			pantry.DELETE("/:itemId", ph.Delete)
			pantry.POST("/:itemId/ack-storage-warning", ph.AckStorageWarning)
		}

		metadata := api.Group("/metadata")
		{
			metadata.GET("/storage-types", mh.ListStorageTypes)
			metadata.GET("/food-categories", mh.ListFoodCategories)
		}

		auth := api.Group("/auth")
		{
			auth.GET("/kakao/login-url", ah.GetKakaoLoginURL)
			auth.GET("/kakao/callback", ah.KakaoCallback)
			auth.POST("/register", ah.Register)
			auth.POST("/refresh", ah.RefreshToken)
			auth.POST("/logout", ah.Logout)
			auth.GET("/session", ah.GetSession)
		}
	}

	swaggerRedirect := func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
	}
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	r.GET("/swagger", swaggerRedirect)
	r.GET("/swagger/*any", func(c *gin.Context) {
		any := c.Param("any")
		if any == "" || any == "/" {
			swaggerRedirect(c)
			return
		}
		swaggerHandler(c)
	})

	return r
}
