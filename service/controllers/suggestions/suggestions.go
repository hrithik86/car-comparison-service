package suggestions

import (
	"car-comparison-service/db/model"
	"context"
	"gorm.io/gorm"
)

type ISuggestionsController interface {
	ExecuteRules(ctx context.Context, db *gorm.DB, vehicle *model.Vehicle) ([]model.VehicleSuggestionResult, error)
}
