package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"naengrijeongdon/apps/server/internal/config"
)

type UserHandler struct {
	cfg  config.Config
	pool *pgxpool.Pool
}

func NewUserHandler(cfg config.Config, pool *pgxpool.Pool) *UserHandler {
	return &UserHandler{cfg: cfg, pool: pool}
}

// GetMe godoc
// @Summary 내 사용자 정보 조회
// @Tags Users
// @Produce json
// @Security AccessTokenAuth
// @Success 200 {object} map[string]any
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	token := strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer "))
	if token == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "authorization bearer token is required"})
		return
	}

	claims := &authClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(h.cfg.AuthJWTSecret), nil
	})
	if err != nil || !parsedToken.Valid || claims.Purpose != tokenPurposeAccess {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "invalid access token"})
		return
	}

	var email string
	var displayName *string
	var provider *string
	var isAdmin bool
	err = h.pool.QueryRow(c.Request.Context(), `
		select email, display_name, provider, is_admin
		from app_users
		where id = $1::uuid
		limit 1
	`, claims.Subject).Scan(&email, &displayName, &provider, &isAdmin)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          claims.Subject,
		"email":       email,
		"displayName": valueOrEmpty(displayName),
		"provider":    valueOrDefaultPtr(provider, "local"),
		"isAdmin":     isAdmin,
	})
}

// UpdateMe godoc
// @Summary 내 사용자 정보 수정
// @Tags Users
// @Security AccessTokenAuth
// @Accept json
// @Produce json
// @Param request body UpdateMeRequest true "수정 요청"
// @Success 200 {object} map[string]any
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /users/me [patch]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	notImplemented(c)
}
