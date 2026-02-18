package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"naengrijeongdon/apps/server/internal/config"
)

const (
	tokenPurposeAccess      = "access"
	tokenPurposeRefresh     = "refresh"
	tokenPurposeSignup      = "signup"
	defaultKakaoRedirectURI = "http://localhost:8080/api/v1/auth/kakao/callback"
)

type AuthHandler struct {
	cfg        config.Config
	pool       *pgxpool.Pool
	httpClient *http.Client
}

type authClaims struct {
	Purpose        string `json:"purpose"`
	Email          string `json:"email,omitempty"`
	Provider       string `json:"provider,omitempty"`
	ProviderUserID string `json:"providerUserId,omitempty"`
	jwt.RegisteredClaims
}

type kakaoTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type kakaoUserResponse struct {
	ID           int64 `json:"id"`
	KakaoAccount struct {
		Email               string `json:"email"`
		HasEmail            bool   `json:"has_email"`
		EmailNeedsAgreement bool   `json:"email_needs_agreement"`
		IsEmailValid        bool   `json:"is_email_valid"`
		IsEmailVerified     bool   `json:"is_email_verified"`
	} `json:"kakao_account"`
}

func NewAuthHandler(cfg config.Config, pool *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{
		cfg:  cfg,
		pool: pool,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetKakaoLoginURL godoc
// @Summary 카카오 OAuth 로그인 URL 조회
// @Description 프론트에서 카카오 인증 화면으로 이동할 URL을 조회합니다.
// @Tags Auth
// @Security AccessTokenAuth
// @Produce json
// @Param redirectUri query string false "로그인 완료 후 리다이렉트 URI" default(http://localhost:8080/api/v1/auth/kakao/callback)
// @Success 200 {object} map[string]any
// @Failure 500 {object} ErrorResponse
// @Router /auth/kakao/login-url [get]
func (h *AuthHandler) GetKakaoLoginURL(c *gin.Context) {
	redirectURI := strings.TrimSpace(c.Query("redirectUri"))
	if redirectURI == "" {
		redirectURI = h.cfg.KakaoRedirectURI
	}
	if redirectURI == "" {
		redirectURI = defaultKakaoRedirectURI
	}

	values := url.Values{}
	values.Set("response_type", "code")
	values.Set("client_id", h.cfg.KakaoRESTAPIKey)
	values.Set("redirect_uri", redirectURI)

	c.JSON(http.StatusOK, gin.H{
		"provider": "kakao",
		"url":      "https://kauth.kakao.com/oauth/authorize?" + values.Encode(),
	})
}

// KakaoCallback godoc
// @Summary 카카오 OAuth 콜백 처리
// @Description 카카오 code를 받아 로그인 또는 회원가입 필요 상태를 반환합니다.
// @Tags Auth
// @Security AccessTokenAuth
// @Produce json
// @Param code query string true "인가 코드"
// @Success 200 {object} map[string]any
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/kakao/callback [get]
func (h *AuthHandler) KakaoCallback(c *gin.Context) {
	code := strings.TrimSpace(c.Query("code"))
	if code == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "code is required"})
		return
	}

	kakaoAccessToken, err := h.exchangeKakaoCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "failed to exchange kakao code"})
		return
	}

	email, providerUserID, err := h.fetchKakaoEmail(c.Request.Context(), kakaoAccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "failed to fetch kakao user info"})
		return
	}
	if email == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "kakao email is required"})
		return
	}

	ctx := c.Request.Context()

	userID, displayName, emailFromProvider, found, err := h.findUserByProviderUserID(ctx, "kakao", providerUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to find user"})
		return
	}
	log.Printf("kakao callback lookup by provider_user_id: provider_user_id=%s found=%t user_id=%s", providerUserID, found, userID)

	if !found {
		userID, displayName, found, err = h.findUserByEmail(ctx, email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to find user"})
			return
		}
		log.Printf("kakao callback lookup by email: email=%s found=%t user_id=%s", email, found, userID)
	}

	if !found {
		signupToken, err := h.issueSignupToken(email, providerUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to issue signup token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"provider":       "kakao",
			"signupRequired": true,
			"email":          email,
			"signupToken":    signupToken,
			"message":        "회원가입이 필요합니다.",
		})
		return
	}

	if strings.TrimSpace(emailFromProvider) != "" {
		email = emailFromProvider
	}

	authPayload, err := h.issueLoginTokens(ctx, userID, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to issue tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"provider":       "kakao",
		"signupRequired": false,
		"user": gin.H{
			"id":          userID,
			"email":       email,
			"displayName": displayName,
			"provider":    "kakao",
		},
		"accessToken":             authPayload.AccessToken,
		"accessTokenExpiresInSec": authPayload.AccessTokenExpiresInSec,
		"refreshToken":            authPayload.RefreshToken,
		"refreshTokenExpiresAt":   authPayload.RefreshTokenExpiresAt,
		"tokenType":               "Bearer",
	})
}

// Register godoc
// @Summary 회원가입
// @Description signupToken 기반으로 회원가입 후 로그인 토큰을 발급합니다.
// @Tags Auth
// @Security AccessTokenAuth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "회원가입 요청"
// @Success 201 {object} map[string]any
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request body"})
		return
	}

	req.DisplayName = strings.TrimSpace(req.DisplayName)
	if req.DisplayName == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "displayName is required"})
		return
	}

	claims, err := h.parseToken(req.SignupToken, tokenPurposeSignup)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid signup token"})
		return
	}

	email := strings.TrimSpace(claims.Email)
	if email == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid signup token payload"})
		return
	}

	ctx := c.Request.Context()
	userID, displayName, found, err := h.findUserByEmail(ctx, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to find user"})
		return
	}

	if !found {
		userID, err = h.createUser(ctx, email, req.DisplayName, claims.Provider, claims.ProviderUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to create user"})
			return
		}
		displayName = req.DisplayName
	}

	authPayload, err := h.issueLoginTokens(ctx, userID, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to issue tokens"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"signupRequired": false,
		"user": gin.H{
			"id":          userID,
			"email":       email,
			"displayName": displayName,
			"provider":    valueOrDefault(claims.Provider, "kakao"),
		},
		"accessToken":             authPayload.AccessToken,
		"accessTokenExpiresInSec": authPayload.AccessTokenExpiresInSec,
		"refreshToken":            authPayload.RefreshToken,
		"refreshTokenExpiresAt":   authPayload.RefreshTokenExpiresAt,
		"tokenType":               "Bearer",
	})
}

// RefreshToken godoc
// @Summary 액세스 토큰 재발급
// @Tags Auth
// @Security AccessTokenAuth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "리프레시 토큰 요청"
// @Security RefreshTokenAuth
// @Success 200 {object} map[string]any
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request body"})
		return
	}

	refreshToken := strings.TrimSpace(valueOrEmpty(req.RefreshToken))
	if refreshToken == "" {
		refreshToken = strings.TrimSpace(strings.TrimPrefix(c.GetHeader("X-Refresh-Token"), "Bearer "))
	}
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "refresh token is required"})
		return
	}

	claims, err := h.parseToken(refreshToken, tokenPurposeRefresh)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "invalid refresh token"})
		return
	}

	tokenHash := hashToken(refreshToken)
	valid, err := h.isRefreshTokenValid(c.Request.Context(), tokenHash, claims.Subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to validate refresh token"})
		return
	}
	if !valid {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "refresh token is revoked or expired"})
		return
	}

	if err := h.revokeRefreshToken(c.Request.Context(), tokenHash); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to revoke refresh token"})
		return
	}

	authPayload, err := h.issueLoginTokens(c.Request.Context(), claims.Subject, claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to issue tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":             authPayload.AccessToken,
		"accessTokenExpiresInSec": authPayload.AccessTokenExpiresInSec,
		"refreshToken":            authPayload.RefreshToken,
		"refreshTokenExpiresAt":   authPayload.RefreshTokenExpiresAt,
		"tokenType":               "Bearer",
	})
}

// Logout godoc
// @Summary 로그아웃
// @Tags Auth
// @Security AccessTokenAuth
// @Accept json
// @Produce json
// @Param request body LogoutRequest false "로그아웃 요청"
// @Security RefreshTokenAuth
// @Success 200 {object} map[string]any
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid request body"})
		return
	}

	refreshToken := strings.TrimSpace(valueOrEmpty(req.RefreshToken))
	if refreshToken == "" {
		refreshToken = strings.TrimSpace(strings.TrimPrefix(c.GetHeader("X-Refresh-Token"), "Bearer "))
	}

	if refreshToken != "" {
		tokenHash := hashToken(refreshToken)
		if err := h.revokeRefreshToken(c.Request.Context(), tokenHash); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to revoke refresh token"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// GetSession godoc
// @Summary 현재 인증 세션 조회
// @Tags Auth
// @Security AccessTokenAuth
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/session [get]
func (h *AuthHandler) GetSession(c *gin.Context) {
	token := strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer "))
	if token == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "authorization bearer token is required"})
		return
	}

	claims, err := h.parseToken(token, tokenPurposeAccess)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "invalid access token"})
		return
	}

	email, displayName, provider, found, err := h.findUserByID(c.Request.Context(), claims.Subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to find user"})
		return
	}
	if !found {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "session not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":          claims.Subject,
			"email":       email,
			"displayName": displayName,
			"provider":    provider,
		},
	})
}

type issuedTokens struct {
	AccessToken             string
	RefreshToken            string
	AccessTokenExpiresInSec int64
	RefreshTokenExpiresAt   string
}

func (h *AuthHandler) issueLoginTokens(ctx context.Context, userID, email string) (issuedTokens, error) {
	now := time.Now()
	accessExp := now.Add(time.Duration(h.cfg.AccessTokenTTLMinutes) * time.Minute)
	refreshExp := now.Add(time.Duration(h.cfg.RefreshTokenTTLDays) * 24 * time.Hour)

	accessClaims := authClaims{
		Purpose: tokenPurposeAccess,
		Email:   email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	accessToken, err := h.signToken(accessClaims)
	if err != nil {
		return issuedTokens{}, err
	}

	refreshClaims := authClaims{
		Purpose: tokenPurposeRefresh,
		Email:   email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshToken, err := h.signToken(refreshClaims)
	if err != nil {
		return issuedTokens{}, err
	}

	if err := h.storeRefreshToken(ctx, userID, refreshToken, refreshExp); err != nil {
		return issuedTokens{}, err
	}

	return issuedTokens{
		AccessToken:             accessToken,
		RefreshToken:            refreshToken,
		AccessTokenExpiresInSec: int64(accessExp.Sub(now).Seconds()),
		RefreshTokenExpiresAt:   refreshExp.UTC().Format(time.RFC3339),
	}, nil
}

func (h *AuthHandler) issueSignupToken(email, providerUserID string) (string, error) {
	now := time.Now()
	claims := authClaims{
		Purpose:        tokenPurposeSignup,
		Email:          email,
		Provider:       "kakao",
		ProviderUserID: providerUserID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(h.cfg.SignupTokenTTLMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	return h.signToken(claims)
}

func (h *AuthHandler) signToken(claims authClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.AuthJWTSecret))
}

func (h *AuthHandler) parseToken(tokenString, requiredPurpose string) (*authClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authClaims{}, func(_ *jwt.Token) (interface{}, error) {
		return []byte(h.cfg.AuthJWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*authClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid claims")
	}

	if claims.Purpose != requiredPurpose {
		return nil, fmt.Errorf("invalid token purpose")
	}

	return claims, nil
}

func (h *AuthHandler) exchangeKakaoCode(ctx context.Context, code string) (string, error) {
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("client_id", h.cfg.KakaoRESTAPIKey)
	values.Set("client_secret", h.cfg.KakaoClientSecret)
	values.Set("redirect_uri", h.cfg.KakaoRedirectURI)
	values.Set("code", code)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://kauth.kakao.com/oauth/token", strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("kakao token status %d", resp.StatusCode)
	}

	var tokenResp kakaoTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}
	if strings.TrimSpace(tokenResp.AccessToken) == "" {
		return "", fmt.Errorf("missing kakao access token")
	}

	return tokenResp.AccessToken, nil
}

func (h *AuthHandler) fetchKakaoEmail(ctx context.Context, kakaoAccessToken string) (string, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+kakaoAccessToken)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	log.Printf("kakao /v2/user/me raw response status=%d body=%s", resp.StatusCode, string(rawBody))

	if resp.StatusCode >= 400 {
		return "", "", fmt.Errorf("kakao user info status %d", resp.StatusCode)
	}

	var userResp kakaoUserResponse
	if err := json.Unmarshal(rawBody, &userResp); err != nil {
		return "", "", err
	}

	if strings.TrimSpace(userResp.KakaoAccount.Email) == "" {
		log.Printf(
			"kakao email missing: user_id=%d has_email=%t needs_agreement=%t is_email_valid=%t is_email_verified=%t",
			userResp.ID,
			userResp.KakaoAccount.HasEmail,
			userResp.KakaoAccount.EmailNeedsAgreement,
			userResp.KakaoAccount.IsEmailValid,
			userResp.KakaoAccount.IsEmailVerified,
		)
	}

	return strings.TrimSpace(userResp.KakaoAccount.Email), fmt.Sprintf("%d", userResp.ID), nil
}

func (h *AuthHandler) findUserByEmail(ctx context.Context, email string) (string, string, bool, error) {
	var id string
	var displayName *string
	err := h.pool.QueryRow(ctx, `
		select id::text, display_name
		from app_users
		where email = $1
		limit 1
	`, email).Scan(&id, &displayName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", false, nil
		}
		return "", "", false, err
	}

	return id, valueOrEmpty(displayName), true, nil
}

func (h *AuthHandler) findUserByProviderUserID(ctx context.Context, provider, providerUserID string) (string, string, string, bool, error) {
	if strings.TrimSpace(providerUserID) == "" {
		return "", "", "", false, nil
	}

	var id string
	var email string
	var displayName *string
	err := h.pool.QueryRow(ctx, `
		select id::text, email, display_name
		from app_users
		where provider = $1
		  and provider_user_id = $2
		limit 1
	`, provider, providerUserID).Scan(&id, &email, &displayName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", "", false, nil
		}
		return "", "", "", false, err
	}

	return id, valueOrEmpty(displayName), email, true, nil
}

func (h *AuthHandler) findUserByID(ctx context.Context, userID string) (string, string, string, bool, error) {
	var email string
	var displayName *string
	var provider *string
	err := h.pool.QueryRow(ctx, `
		select email, display_name, provider
		from app_users
		where id = $1::uuid
		limit 1
	`, userID).Scan(&email, &displayName, &provider)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", "", false, nil
		}
		return "", "", "", false, err
	}
	return email, valueOrEmpty(displayName), valueOrDefaultPtr(provider, "local"), true, nil
}

func (h *AuthHandler) createUser(ctx context.Context, email, displayName, provider, providerUserID string) (string, error) {
	var id string
	provider = valueOrDefault(provider, "kakao")
	providerUserID = strings.TrimSpace(providerUserID)

	err := h.pool.QueryRow(ctx, `
		insert into app_users (email, display_name, provider, provider_user_id)
		values ($1, $2, $3, nullif($4, ''))
		returning id::text
	`, email, displayName, provider, providerUserID).Scan(&id)
	return id, err
}

func (h *AuthHandler) storeRefreshToken(ctx context.Context, userID, refreshToken string, expiresAt time.Time) error {
	_, err := h.pool.Exec(ctx, `
		insert into auth_refresh_tokens (user_id, token_hash, expires_at)
		values ($1::uuid, $2, $3)
	`, userID, hashToken(refreshToken), expiresAt)
	return err
}

func (h *AuthHandler) isRefreshTokenValid(ctx context.Context, tokenHash, userID string) (bool, error) {
	var exists bool
	err := h.pool.QueryRow(ctx, `
		select exists (
			select 1
			from auth_refresh_tokens
			where token_hash = $1
			  and user_id = $2::uuid
			  and revoked_at is null
			  and expires_at > now()
		)
	`, tokenHash, userID).Scan(&exists)
	return exists, err
}

func (h *AuthHandler) revokeRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := h.pool.Exec(ctx, `
		update auth_refresh_tokens
		set revoked_at = now()
		where token_hash = $1
		  and revoked_at is null
	`, tokenHash)
	return err
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func valueOrDefault(s, fallback string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return fallback
	}
	return s
}

func valueOrDefaultPtr(s *string, fallback string) string {
	if s == nil {
		return fallback
	}
	return valueOrDefault(*s, fallback)
}
