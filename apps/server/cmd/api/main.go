package main

import (
	"log"

	"naengrijeongdon/apps/server/internal/config"
	"naengrijeongdon/apps/server/internal/db"
	"naengrijeongdon/apps/server/internal/router"
)

// @title NaengriJeongdon API
// @version 0.1.0
// @description API for naengrijeongdon inventory management.
// @BasePath /api/v1
// @securityDefinitions.apikey AccessTokenAuth
// @in header
// @name Authorization
// @description Bearer access token. Example: Bearer eyJ...
// @securityDefinitions.apikey RefreshTokenAuth
// @in header
// @name X-Refresh-Token
// @description Bearer refresh token. Example: Bearer eyJ...
// @security AccessTokenAuth
// @security RefreshTokenAuth
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	pool, err := db.NewPostgresPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer pool.Close()

	r := router.New(cfg, pool)

	log.Printf("server listening on %s", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
