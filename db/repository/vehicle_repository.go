package repository

import (
	"car-comparison-service/db/model"
	"car-comparison-service/db/utils"
	"context"
	"github.com/google/uuid"
)

type IVehicle interface {
	GetVehiclesByModel(ctx context.Context, vehicleName string) ([]*model.Vehicle, error)
	GetVehiclesById(ctx context.Context, id uuid.UUID) (*model.Vehicle, error)
	GetVehiclesByIds(ctx context.Context, ids []uuid.UUID) ([]*model.Vehicle, error)
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

func (db CarComparisonServiceDb) GetVehiclesById(ctx context.Context, id uuid.UUID) (*model.Vehicle, error) {
	var vehicle *model.Vehicle
	result := db.WithContext(ctx).
		Table(model.TableNameVehicle).
		Where("id = ?", id).
		Take(&vehicle)
	err := utils.ValidateResultSuccess(result)
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (db CarComparisonServiceDb) GetVehiclesByIds(ctx context.Context, ids []uuid.UUID) ([]*model.Vehicle, error) {
	var vehicles []*model.Vehicle
	result := db.WithContext(ctx).
		Table(model.TableNameVehicle).
		Where("id in ?", ids).
		Scan(&vehicles)
	err := utils.ValidateResultSuccess(result)
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}
