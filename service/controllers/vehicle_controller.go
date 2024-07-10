package controllers

import (
	"car-comparison-service/db/model"
	"car-comparison-service/db/repository"
	"context"
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
