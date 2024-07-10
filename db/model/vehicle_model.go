package model

import "github.com/google/uuid"

const TableNameVehicle = "vehicle"

type Vehicle struct {
	DbId
	Model             *string `gorm:"column:model;not null" json:"model"`
	Brand             *string `gorm:"column:brand;not null" json:"brand"`
	ManufacturingYear *int    `gorm:"column:manufacturing_year;not null" json:"manufacturing_year"`
	DbTimeAudit
}

func (vehicle *Vehicle) GetPrimaryID() uuid.UUID {
	return *vehicle.ID
}

func (vehicle *Vehicle) GetTableName() string {
	return TableNameVehicle
}
