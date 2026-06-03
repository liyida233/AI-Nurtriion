package database

import (
	"fmt"
	"time"

	"ai-nutrition/backend/internal/config"
	"ai-nutrition/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.Exercise{},
		&models.WorkoutSession{},
		&models.WorkoutEntry{},
		&models.FoodItem{},
		&models.MealLog{},
		&models.BodyRecord{},
		&models.Goal{},
		&models.GoalMilestone{},
		&models.AnalyticsSnapshot{},
		&models.AIRecommendation{},
		&models.RecommendationFeedback{},
		&models.ProgressReport{},
	)
}

func Seed(db *gorm.DB) error {
	exercises := []models.Exercise{
		{Name: "Squat", Category: "strength", MuscleGroup: "lower_body", Equipment: "barbell", IntensityLevel: "moderate"},
		{Name: "Bench Press", Category: "strength", MuscleGroup: "chest", Equipment: "barbell", IntensityLevel: "moderate"},
		{Name: "Lat Pulldown", Category: "strength", MuscleGroup: "back", Equipment: "machine", IntensityLevel: "moderate"},
		{Name: "Running", Category: "cardio", MuscleGroup: "full_body", Equipment: "none", IntensityLevel: "variable"},
	}

	for _, exercise := range exercises {
		var existing models.Exercise
		err := db.Where("name = ?", exercise.Name).First(&existing).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		exercise.ID = uuid.NewString()
		if err := db.Create(&exercise).Error; err != nil {
			return err
		}
	}

	foods := []models.FoodItem{
		{Name: "Chicken Breast", ServingSize: "100g", Calories: 165, Protein: 31, Carbohydrates: 0, Fat: 3.6},
		{Name: "Cooked Rice", ServingSize: "100g", Calories: 130, Protein: 2.7, Carbohydrates: 28, Fat: 0.3},
		{Name: "Egg", ServingSize: "1 large", Calories: 72, Protein: 6.3, Carbohydrates: 0.4, Fat: 4.8},
		{Name: "Banana", ServingSize: "1 medium", Calories: 105, Protein: 1.3, Carbohydrates: 27, Fat: 0.4},
	}

	for _, food := range foods {
		var existing models.FoodItem
		err := db.Where("name = ?", food.Name).First(&existing).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		food.ID = uuid.NewString()
		if err := db.Create(&food).Error; err != nil {
			return err
		}
	}

	return nil
}
