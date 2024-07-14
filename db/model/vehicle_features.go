package model

import "github.com/google/uuid"

const TableNameVehicleFeatures = "vehicle_features"

type VehicleFeatures struct {
	DbId
	Key       *string    `gorm:"column:key" json:"key"`
	Value     *string    `gorm:"column:value" json:"value"`
	VehicleId *uuid.UUID `gorm:"column:vehicle_id" json:"vehicle_id"`
	DbTimeAudit
}
