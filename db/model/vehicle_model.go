package model

import "github.com/google/uuid"

const TableNameVehicle = "vehicle"

type VehicleType string

const (
	CAR   VehicleType = "CAR"
	TRUCK VehicleType = "TRUCK"
)

type Vehicle struct {
	DbId
	Model             *string      `gorm:"column:model;not null" json:"model"`
	Brand             *string      `gorm:"column:brand;not null" json:"brand"`
	ManufacturingYear *int         `gorm:"column:manufacturing_year;not null" json:"manufacturing_year"`
	Type              *VehicleType `gorm:"column:type;not null" json:"vehicle_type"`
	Price             *int64       `gorm:"column:price;not null" json:"price"`
	DbTimeAudit
}

type VehicleWithAttachmentInformation struct {
	Id                *uuid.UUID   `gorm:"column:id"`
	Model             *string      `gorm:"column:model"`
	Brand             *string      `gorm:"column:brand"`
	ManufacturingYear *int         `gorm:"column:manufacturing_year"`
	Type              *VehicleType `gorm:"column:type"`
	Price             *int64       `gorm:"column:price"`
	AttachmentId      *uuid.UUID   `gorm:"column:attachment_id"`
	Path              *string      `gorm:"column:path"`
	MediaType         *string      `gorm:"column:media_type"`
}
