package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"naengrijeongdon/apps/server/internal/catalog"
)

type MetadataHandler struct{}

func NewMetadataHandler() *MetadataHandler { return &MetadataHandler{} }

// ListStorageTypes godoc
// @Summary 보관 타입 목록 조회
// @Description 코드 설명: FRIDGE(냉장), FREEZER(냉동), KIMCHI(김치냉장), ROOM(실온)
// @Tags Metadata
// @Security AccessTokenAuth
// @Produce json
// @Success 200 {array} catalog.CodeName
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /metadata/storage-types [get]
func (h *MetadataHandler) ListStorageTypes(c *gin.Context) {
	c.JSON(http.StatusOK, catalog.StorageTypeCatalog)
}

// ListFoodCategories godoc
// @Summary 식품 카테고리 목록 조회
// @Description 코드 설명: DAIRY(유제품), MEAT(육류), SEAFOOD(해산물), VEGETABLE(채소), FRUIT(과일), KIMCHI_SIDE(김치/반찬), FROZEN_FOOD(냉동식품), SAUCE(소스/조미료), DRINK(음료), ETC(기타)
// @Tags Metadata
// @Security AccessTokenAuth
// @Produce json
// @Success 200 {array} catalog.CodeName
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /metadata/food-categories [get]
func (h *MetadataHandler) ListFoodCategories(c *gin.Context) {
	c.JSON(http.StatusOK, catalog.FoodCategoryCatalog)
}
