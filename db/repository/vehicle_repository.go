package repository

import (
	"car-comparison-service/db/model"
	"car-comparison-service/db/utils"
	"context"
)

type IVehicle interface {
	GetVehiclesByModel(ctx context.Context, vehicleName string) ([]*model.Vehicle, error)
}

func (db CarComparisonServiceDb) GetVehiclesByModel(ctx context.Context, modelName string) ([]*model.Vehicle, error) {
	var vehicles []*model.Vehicle
	result := db.WithContext(ctx).
		Table(model.TableNameVehicle).
		Where("model like ?", "%"+modelName+"%").
		Find(&vehicles)
	err := utils.ValidateResultSuccess(result)
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}
