package handlers

import "github.com/gin-gonic/gin"

type FoodHandler struct{}

func NewFoodHandler() *FoodHandler { return &FoodHandler{} }

// AutocompleteFoods godoc
// @Summary 식품 자동완성 조회
// @Tags Foods
// @Security AccessTokenAuth
// @Produce json
// @Param q query string true "검색어" example(우유)
// @Param limit query int false "최대 개수" default(10) minimum(1) maximum(50)
// @Success 200 {object} map[string]any
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /foods/autocomplete [get]
func (h *FoodHandler) Autocomplete(c *gin.Context) {
	notImplemented(c)
}

// GetFood godoc
// @Summary 식품 상세 조회
// @Tags Foods
// @Security AccessTokenAuth
// @Produce json
// @Param foodId path string true "식품 ID"
// @Success 200 {object} map[string]any
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /foods/{foodId} [get]
func (h *FoodHandler) Get(c *gin.Context) {
	notImplemented(c)
}

// PredictExpiry godoc
// @Summary 유통기한 예측 미리보기
// @Tags Foods
// @Security AccessTokenAuth
// @Produce json
// @Param foodId path string true "식품 ID"
// @Param storageSpaceId query string true "보관 공간 ID" example(00000000-0000-0000-0000-000000000005)
// @Param purchasedAt query string true "구매일(YYYY-MM-DD)" example(2026-02-18)
// @Success 200 {object} map[string]any
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /foods/{foodId}/predict-expiry [get]
func (h *FoodHandler) PredictExpiry(c *gin.Context) {
	notImplemented(c)
}
