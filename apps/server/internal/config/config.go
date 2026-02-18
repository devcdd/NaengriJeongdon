package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort              string
	DatabaseURL           string
	KakaoRESTAPIKey       string
	KakaoClientSecret     string
	KakaoRedirectURI      string
	AuthJWTSecret         string
	AccessTokenTTLMinutes int
	RefreshTokenTTLDays   int
	SignupTokenTTLMinutes int
}

func Load() (Config, error) {
	_ = godotenv.Load()

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	kakaoRESTAPIKey := os.Getenv("KAKAO_REST_API_KEY")
	if kakaoRESTAPIKey == "" {
		return Config{}, fmt.Errorf("KAKAO_REST_API_KEY is required")
	}

	kakaoClientSecret := os.Getenv("KAKAO_CLIENT_SECRET")
	if kakaoClientSecret == "" {
		return Config{}, fmt.Errorf("KAKAO_CLIENT_SECRET is required")
	}

	kakaoRedirectURI := os.Getenv("KAKAO_REDIRECT_URI")
	if kakaoRedirectURI == "" {
		return Config{}, fmt.Errorf("KAKAO_REDIRECT_URI is required")
	}

	authJWTSecret := os.Getenv("AUTH_JWT_SECRET")
	if authJWTSecret == "" {
		return Config{}, fmt.Errorf("AUTH_JWT_SECRET is required")
	}

	accessTokenTTLMinutes := 60
	if os.Getenv("AUTH_ACCESS_TOKEN_TTL_MINUTES") != "" {
		_, err := fmt.Sscanf(os.Getenv("AUTH_ACCESS_TOKEN_TTL_MINUTES"), "%d", &accessTokenTTLMinutes)
		if err != nil || accessTokenTTLMinutes <= 0 {
			return Config{}, fmt.Errorf("AUTH_ACCESS_TOKEN_TTL_MINUTES must be a positive integer")
		}
	}

	refreshTokenTTLDays := 30
	if os.Getenv("AUTH_REFRESH_TOKEN_TTL_DAYS") != "" {
		_, err := fmt.Sscanf(os.Getenv("AUTH_REFRESH_TOKEN_TTL_DAYS"), "%d", &refreshTokenTTLDays)
		if err != nil || refreshTokenTTLDays <= 0 {
			return Config{}, fmt.Errorf("AUTH_REFRESH_TOKEN_TTL_DAYS must be a positive integer")
		}
	}

	signupTokenTTLMinutes := 30
	if os.Getenv("AUTH_SIGNUP_TOKEN_TTL_MINUTES") != "" {
		_, err := fmt.Sscanf(os.Getenv("AUTH_SIGNUP_TOKEN_TTL_MINUTES"), "%d", &signupTokenTTLMinutes)
		if err != nil || signupTokenTTLMinutes <= 0 {
			return Config{}, fmt.Errorf("AUTH_SIGNUP_TOKEN_TTL_MINUTES must be a positive integer")
		}
	}

	return Config{
		HTTPPort:              port,
		DatabaseURL:           databaseURL,
		KakaoRESTAPIKey:       kakaoRESTAPIKey,
		KakaoClientSecret:     kakaoClientSecret,
		KakaoRedirectURI:      kakaoRedirectURI,
		AuthJWTSecret:         authJWTSecret,
		AccessTokenTTLMinutes: accessTokenTTLMinutes,
		RefreshTokenTTLDays:   refreshTokenTTLDays,
		SignupTokenTTLMinutes: signupTokenTTLMinutes,
	}, nil
}
