package catalog

type StorageTypeCode string

const (
	StorageTypeFridge  StorageTypeCode = "FRIDGE"
	StorageTypeFreezer StorageTypeCode = "FREEZER"
	StorageTypeKimchi  StorageTypeCode = "KIMCHI"
	StorageTypeRoom    StorageTypeCode = "ROOM"
)

type FoodCategoryCode string

const (
	FoodCategoryDairy      FoodCategoryCode = "DAIRY"
	FoodCategoryMeat       FoodCategoryCode = "MEAT"
	FoodCategorySeafood    FoodCategoryCode = "SEAFOOD"
	FoodCategoryVegetable  FoodCategoryCode = "VEGETABLE"
	FoodCategoryFruit      FoodCategoryCode = "FRUIT"
	FoodCategoryKimchiSide FoodCategoryCode = "KIMCHI_SIDE"
	FoodCategoryFrozenFood FoodCategoryCode = "FROZEN_FOOD"
	FoodCategorySauce      FoodCategoryCode = "SAUCE"
	FoodCategoryDrink      FoodCategoryCode = "DRINK"
	FoodCategoryEtc        FoodCategoryCode = "ETC"
)

type CodeName struct {
	Code string `json:"code" example:"FRIDGE"`
	Name string `json:"name" example:"냉장"`
}

var StorageTypeCatalog = []CodeName{
	{Code: string(StorageTypeFridge), Name: "냉장"},
	{Code: string(StorageTypeFreezer), Name: "냉동"},
	{Code: string(StorageTypeKimchi), Name: "김치냉장"},
	{Code: string(StorageTypeRoom), Name: "실온"},
}

var FoodCategoryCatalog = []CodeName{
	{Code: string(FoodCategoryDairy), Name: "유제품"},
	{Code: string(FoodCategoryMeat), Name: "육류"},
	{Code: string(FoodCategorySeafood), Name: "해산물"},
	{Code: string(FoodCategoryVegetable), Name: "채소"},
	{Code: string(FoodCategoryFruit), Name: "과일"},
	{Code: string(FoodCategoryKimchiSide), Name: "김치/반찬"},
	{Code: string(FoodCategoryFrozenFood), Name: "냉동식품"},
	{Code: string(FoodCategorySauce), Name: "소스/조미료"},
	{Code: string(FoodCategoryDrink), Name: "음료"},
	{Code: string(FoodCategoryEtc), Name: "기타"},
}
