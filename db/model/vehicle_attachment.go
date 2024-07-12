package model

import "github.com/google/uuid"

const TableNameVehicleAttachment = "vehicle_attachment"

type MediaType string

const (
	IMAGE MediaType = "IMAGE"
	VIDEO MediaType = "VIDEO"
)

type VehicleAttachment struct {
	DbId
	Name      *string    `gorm:"column:name" json:"name"`
	Path      *string    `gorm:"column:path" json:"path"`
	MediaType *MediaType `gorm:"column:mediatype" json:"mediatype"`
	VehicleId *uuid.UUID `gorm:"column:vehicle_id" json:"vehicle_id"`
	DbTimeAudit
}
