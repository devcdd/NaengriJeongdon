package handlers

import "github.com/gin-gonic/gin"

type PantryItemHandler struct{}

func NewPantryItemHandler() *PantryItemHandler { return &PantryItemHandler{} }

// ListPantryItems godoc
// @Summary 내 식재료 목록 조회
// @Description status 코드: ALL(전체), EXPIRED(만료), TODAY(오늘 마감), DUE_3_DAYS(3일 이내), DUE_7_DAYS(7일 이내), SAFE(여유 있음), NO_EXPIRY(마감일 없음). sort 코드: EXPIRES_ASC(마감일 빠른순), CREATED_DESC(등록 최신순).
// @Tags PantryItems
// @Security AccessTokenAuth
// @Produce json
// @Param status query string false "상태 필터" Enums(ALL,EXPIRED,TODAY,DUE_3_DAYS,DUE_7_DAYS,SAFE,NO_EXPIRY) default(ALL)
// @Param storageSpaceId query string false "공간 ID" example(00000000-0000-0000-0000-000000000005)
// @Param q query string false "검색어" example(우유)
// @Param sort query string false "정렬" Enums(EXPIRES_ASC,CREATED_DESC) default(EXPIRES_ASC)
// @Success 200 {object} map[string]any
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /pantry-items [get]
func (h *PantryItemHandler) List(c *gin.Context) {
	notImplemented(c)
}

// CreatePantryItem godoc
// @Summary 식재료 등록
// @Description foodCategoryCode 코드: DAIRY(유제품), MEAT(육류), SEAFOOD(해산물), VEGETABLE(채소), FRUIT(과일), KIMCHI_SIDE(김치/반찬), FROZEN_FOOD(냉동식품), SAUCE(소스/조미료), DRINK(음료), ETC(기타). 실시간 목록은 GET /api/v1/metadata/food-categories 참고.
// @Tags PantryItems
// @Security AccessTokenAuth
// @Accept json
// @Produce json
// @Param request body CreatePantryItemRequest true "등록 요청"
// @Success 201 {object} map[string]any
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /pantry-items [post]
func (h *PantryItemHandler) Create(c *gin.Context) {
	notImplemented(c)
}

// GetPantryItem godoc
// @Summary 식재료 단건 조회
// @Tags PantryItems
// @Security AccessTokenAuth
// @Produce json
// @Param itemId path string true "식재료 ID"
// @Success 200 {object} map[string]any
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /pantry-items/{itemId} [get]
func (h *PantryItemHandler) Get(c *gin.Context) {
	notImplemented(c)
}

// UpdatePantryItem godoc
// @Summary 식재료 수정
// @Tags PantryItems
// @Security AccessTokenAuth
// @Accept json
// @Produce json
// @Param itemId path string true "식재료 ID"
// @Param request body UpdatePantryItemRequest true "수정 요청"
// @Success 200 {object} map[string]any
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /pantry-items/{itemId} [patch]
func (h *PantryItemHandler) Update(c *gin.Context) {
	notImplemented(c)
}

// DeletePantryItem godoc
// @Summary 식재료 삭제
// @Tags PantryItems
// @Security AccessTokenAuth
// @Produce json
// @Param itemId path string true "식재료 ID"
// @Success 204 {string} string "No Content"
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /pantry-items/{itemId} [delete]
func (h *PantryItemHandler) Delete(c *gin.Context) {
	notImplemented(c)
}

// AckStorageWarning godoc
// @Summary 보관 방식 경고 확인 처리
// @Tags PantryItems
// @Security AccessTokenAuth
// @Produce json
// @Param itemId path string true "식재료 ID"
// @Success 200 {object} map[string]any
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /pantry-items/{itemId}/ack-storage-warning [post]
func (h *PantryItemHandler) AckStorageWarning(c *gin.Context) {
	notImplemented(c)
}
