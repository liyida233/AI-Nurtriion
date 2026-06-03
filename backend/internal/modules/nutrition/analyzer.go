package nutrition

import "ai-nutrition/backend/internal/models"

func CalculateMealTotals(food models.FoodItem, quantity float64) MealIndicators {
	return MealIndicators{
		Calories: food.Calories * quantity,
		Protein:  food.Protein * quantity,
		Carbs:    food.Carbohydrates * quantity,
		Fat:      food.Fat * quantity,
	}
}
