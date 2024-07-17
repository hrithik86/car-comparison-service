package repository

import (
	"car-comparison-service/db/model"
	"car-comparison-service/db/utils"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IVehicle interface {
	GetVehiclesByModel(ctx context.Context, vehicleName string) ([]*model.VehicleWithAttachmentInformation, error)
	GetVehicleInfoById(ctx context.Context, id uuid.UUID) (*model.Vehicle, error)
	GetVehiclesByIds(ctx context.Context, ids []uuid.UUID) ([]*model.Vehicle, error)
	GetVehicleWithFeaturesById(ctx context.Context, id uuid.UUID) ([]*model.VehicleWithFeatures, error)
	CreateVehicle(ctx context.Context, vehicle *model.Vehicle) (*model.Vehicle, error)
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

func (db CarComparisonServiceDb) GetVehicleInfoById(ctx context.Context, id uuid.UUID) (*model.Vehicle, error) {
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
		Find(&vehicle)
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
		Find(&vehicles)
	err := utils.ValidateResultSuccess(result)
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (db CarComparisonServiceDb) CreateVehicle(ctx context.Context, vehicle *model.Vehicle) (*model.Vehicle, error) {
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&vehicle).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, utils.ValidateError(err)
	}
	return vehicle, err
}
