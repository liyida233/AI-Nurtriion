package nutrition

type FoodRequest struct {
	Name          string  `json:"name" binding:"required"`
	ServingSize   string  `json:"servingSize"`
	Calories      float64 `json:"calories" binding:"min=0"`
	Protein       float64 `json:"protein" binding:"min=0"`
	Carbohydrates float64 `json:"carbohydrates" binding:"min=0"`
	Fat           float64 `json:"fat" binding:"min=0"`
	Sugar         float64 `json:"sugar" binding:"min=0"`
	Sodium        float64 `json:"sodium" binding:"min=0"`
}

type MealRequest struct {
	FoodItemID string  `json:"foodItemId" binding:"required"`
	MealType   string  `json:"mealType" binding:"required"`
	Quantity   float64 `json:"quantity" binding:"required,min=0.01"`
	MealTime   string  `json:"mealTime" binding:"required"`
}

type MealIndicators struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
}
