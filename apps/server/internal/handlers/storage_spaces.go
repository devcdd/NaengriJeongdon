package handlers

import "github.com/gin-gonic/gin"

type StorageSpaceHandler struct{}

func NewStorageSpaceHandler() *StorageSpaceHandler { return &StorageSpaceHandler{} }

// ListStorageSpaces godoc
// @Summary 보관 공간 목록 조회
// @Tags StorageSpaces
// @Security AccessTokenAuth
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /storage-spaces [get]
func (h *StorageSpaceHandler) List(c *gin.Context) {
	notImplemented(c)
}

// CreateStorageSpace godoc
// @Summary 보관 공간 생성
// @Description storageTypeCode 코드: FRIDGE(냉장), FREEZER(냉동), KIMCHI(김치냉장), ROOM(실온). 실시간 목록은 GET /api/v1/metadata/storage-types 참고.
// @Tags StorageSpaces
// @Security AccessTokenAuth
// @Accept json
// @Produce json
// @Param request body CreateStorageSpaceRequest true "생성 요청"
// @Success 201 {object} map[string]any
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /storage-spaces [post]
func (h *StorageSpaceHandler) Create(c *gin.Context) {
	notImplemented(c)
}

// GetStorageSpace godoc
// @Summary 보관 공간 단건 조회
// @Tags StorageSpaces
// @Security AccessTokenAuth
// @Produce json
// @Param spaceId path string true "공간 ID"
// @Success 200 {object} map[string]any
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /storage-spaces/{spaceId} [get]
func (h *StorageSpaceHandler) Get(c *gin.Context) {
	notImplemented(c)
}

// UpdateStorageSpace godoc
// @Summary 보관 공간 수정
// @Tags StorageSpaces
// @Security AccessTokenAuth
// @Accept json
// @Produce json
// @Param spaceId path string true "공간 ID"
// @Param request body UpdateStorageSpaceRequest true "수정 요청"
// @Success 200 {object} map[string]any
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /storage-spaces/{spaceId} [patch]
func (h *StorageSpaceHandler) Update(c *gin.Context) {
	notImplemented(c)
}

// DeleteStorageSpace godoc
// @Summary 보관 공간 삭제
// @Tags StorageSpaces
// @Security AccessTokenAuth
// @Produce json
// @Param spaceId path string true "공간 ID"
// @Success 204 {string} string "No Content"
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 501 {object} ErrorResponse
// @Router /storage-spaces/{spaceId} [delete]
func (h *StorageSpaceHandler) Delete(c *gin.Context) {
	notImplemented(c)
}
