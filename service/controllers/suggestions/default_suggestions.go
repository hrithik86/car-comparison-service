package suggestions

import (
	"car-comparison-service/db/model"
	"car-comparison-service/ruleEngine/rules/suggestions"
	"context"
	"gorm.io/gorm"
)

type DefaultSuggestionsController struct{}

func NewDefaultSuggestionsController() ISuggestionsController {
	return &DefaultSuggestionsController{}
}

func (c *DefaultSuggestionsController) ExecuteRules(ctx context.Context, db *gorm.DB, vehicle *model.Vehicle) ([]model.VehicleSuggestionResult, error) {
	return suggestions.ExecuteRules(ctx, db, vehicle)
}
