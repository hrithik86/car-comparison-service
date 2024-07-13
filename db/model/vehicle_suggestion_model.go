package model

import (
	"github.com/google/uuid"
)

type VehicleSuggestionResult struct {
	Id                uuid.UUID `json:"id"`
	Model             string    `gorm:"column:model;not null" json:"model"`
	Brand             string    `gorm:"column:brand;not null" json:"brand"`
	ManufacturingYear int       `gorm:"column:manufacturing_year;not null" json:"manufacturing_year"`
	Price             int64     `gorm:"column:price;not null" json:"price"`
	Mileage           float64   `gorm:"column:mileage;not null" json:"mileage"`
	FuelType          string    `gorm:"column:fuel_type;not null" json:"fuel_type"`
	Type              string    `gorm:"column:type;not null" json:"type"`
	Rank              int64     `json:"rank"`
}
