package controllers

import (
	"car-comparison-service/db/model"
	"car-comparison-service/db/repository"
	"car-comparison-service/ruleEngine/rules/suggestions"
	"car-comparison-service/service/api/request"
	"context"
	"github.com/google/uuid"
)

type VehicleController struct {
	db repository.CarComparisonServiceDb
}

func NewVehicleController() *VehicleController {
	return &VehicleController{db: repository.DbClient()}
}

func (vc *VehicleController) GetVehiclesByModelName(ctx context.Context, modelName string) ([]*model.Vehicle, error) {
	vehicles, err := vc.db.GetVehiclesByModel(ctx, modelName)
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (vc *VehicleController) GetVehicleById(ctx context.Context, id uuid.UUID) (*model.Vehicle, error) {
	vehicle, err := vc.db.GetVehiclesById(ctx, id)
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (vc *VehicleController) GetVehicleSuggestions(ctx context.Context, id uuid.UUID) ([]model.VehicleSuggestionResult, error) {
	vehicle, err := vc.db.GetVehiclesById(ctx, id)
	if err != nil {
		return nil, err
	}
	suggestedVehicles, err := suggestions.ExecuteRules(ctx, vc.db.DB, vehicle)
	if err != nil {
		return nil, err
	}
	return suggestedVehicles, nil
}

func (vc *VehicleController) GetVehicleComparison(ctx context.Context, req request.VehicleComparisonRequest) ([]*model.Vehicle, error) {
	vehicles, err := vc.db.GetVehiclesByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}
