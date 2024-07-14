package controllers

import (
	"car-comparison-service/db/model"
	"car-comparison-service/ruleEngine/rules/suggestions"
	"context"
	"gorm.io/gorm"
)

type TruckSuggestionsController struct{}

func NewTruckSuggestionsController() ISuggestionsController {
	return &TruckSuggestionsController{}
}

func (c *TruckSuggestionsController) ExecuteRules(ctx context.Context, db *gorm.DB, vehicle *model.Vehicle) ([]model.VehicleSuggestionResult, error) {
	return suggestions.ExecuteRules(ctx, db, vehicle)
}
