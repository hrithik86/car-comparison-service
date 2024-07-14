package controllers

import (
	"car-comparison-service/db/model"
	"car-comparison-service/ruleEngine/rules/suggestions"
	"context"
	"gorm.io/gorm"
)

type CarSuggestionsController struct{}

func NewCarSuggestionsController() ISuggestionsController {
	return &CarSuggestionsController{}
}

func (c *CarSuggestionsController) ExecuteRules(ctx context.Context, db *gorm.DB, vehicle *model.Vehicle) ([]model.VehicleSuggestionResult, error) {
	return suggestions.ExecuteRules(ctx, db, vehicle)
}
