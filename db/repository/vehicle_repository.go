package repository

import (
	"car-comparison-service/db/model"
	"car-comparison-service/db/utils"
	"context"
	"github.com/google/uuid"
)

type IVehicle interface {
	GetVehiclesByModel(ctx context.Context, vehicleName string) ([]*model.VehicleWithAttachmentInformation, error)
	GetVehiclesById(ctx context.Context, id uuid.UUID) (*model.Vehicle, error)
	GetVehiclesByIds(ctx context.Context, ids []uuid.UUID) ([]*model.Vehicle, error)
	GetVehicleWithFeaturesById(ctx context.Context, id uuid.UUID) (*model.VehicleWithFeatures, error)
}

func (db CarComparisonServiceDb) GetVehiclesByModel(ctx context.Context, modelName string) ([]*model.VehicleWithAttachmentInformation, error) {
	var vehicles []*model.VehicleWithAttachmentInformation
	result := db.WithContext(ctx).
		Table(model.TableNameVehicle).
		Select("vehicle.*, vehicle_attachment.path as path, vehicle_attachment.media_type as media_type, "+
			"vehicle_attachment.id as attachment_id").
		Joins("left join vehicle_attachment on vehicle_attachment.vehicle_id = vehicle.id").
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

func (db CarComparisonServiceDb) GetVehicleWithFeaturesById(ctx context.Context, id uuid.UUID) ([]*model.VehicleWithFeatures, error) {
	var vehicle []*model.VehicleWithFeatures
	result := db.WithContext(ctx).
		Table(model.TableNameVehicle).
		Select("vehicle.*, vehicle_features.id as feature_id, vehicle_features.key as key, vehicle_features.value as value").
		Joins("left join vehicle_features on vehicle_features.vehicle_id = vehicle.id").
		Where("vehicle.id = ?", id).
		Scan(&vehicle)
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
