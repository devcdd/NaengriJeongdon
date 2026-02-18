package handlers

import "naengrijeongdon/apps/server/internal/catalog"

type UpdateMeRequest struct {
	DisplayName string `json:"displayName" example:"냉리정돈고수"`
}

type CreateStorageSpaceRequest struct {
	Name string `json:"name" example:"우리집 메인 냉장고"`
	// 보관 공간 타입 코드: FRIDGE(냉장), FREEZER(냉동), KIMCHI(김치냉장), ROOM(실온)
	StorageTypeCode catalog.StorageTypeCode `json:"storageTypeCode" enums:"FRIDGE,FREEZER,KIMCHI,ROOM" example:"FRIDGE"`
}

type UpdateStorageSpaceRequest struct {
	Name       *string `json:"name,omitempty" example:"작은 냉동실"`
	IsArchived *bool   `json:"isArchived,omitempty" default:"false"`
}

type CreatePantryItemRequest struct {
	StorageSpaceID string  `json:"storageSpaceId" example:"00000000-0000-0000-0000-000000000005"`
	FoodID         *string `json:"foodId,omitempty" example:"00000000-0000-0000-0000-000000000002"`
	// 식품 카테고리 코드: DAIRY, MEAT, SEAFOOD, VEGETABLE, FRUIT, KIMCHI_SIDE, FROZEN_FOOD, SAUCE, DRINK, ETC
	FoodCategory *catalog.FoodCategoryCode `json:"foodCategoryCode,omitempty" enums:"DAIRY,MEAT,SEAFOOD,VEGETABLE,FRUIT,KIMCHI_SIDE,FROZEN_FOOD,SAUCE,DRINK,ETC" example:"DAIRY"`
	InputName    string                    `json:"inputName" example:"우유"`
	PurchasedAt  *string                   `json:"purchasedAt,omitempty" example:"2026-02-14"`
	ExpiresAt    *string                   `json:"expiresAt,omitempty" example:"2026-02-21"`
}

type UpdatePantryItemRequest struct {
	StorageSpaceID *string `json:"storageSpaceId,omitempty" example:"00000000-0000-0000-0000-000000000005"`
	InputName      *string `json:"inputName,omitempty" example:"딸기 우유"`
	PurchasedAt    *string `json:"purchasedAt,omitempty" example:"2026-02-15"`
	ExpiresAt      *string `json:"expiresAt,omitempty" example:"2026-02-22"`
	Notes          *string `json:"notes,omitempty" example:"아침 시리얼용"`
}

type RefreshTokenRequest struct {
	RefreshToken *string `json:"refreshToken,omitempty" example:"sample-refresh-token"`
}

type LogoutRequest struct {
	RefreshToken *string `json:"refreshToken,omitempty" example:"sample-refresh-token"`
}

type RegisterRequest struct {
	SignupToken string `json:"signupToken" example:"sample-signup-token"`
	DisplayName string `json:"displayName" example:"냉장고고수"`
}
