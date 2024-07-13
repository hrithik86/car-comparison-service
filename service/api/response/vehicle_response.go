package response

import "github.com/google/uuid"

type VehicleResponse struct {
	Id                *uuid.UUID          `json:"id"`
	Brand             *string             `json:"brand"`
	Model             *string             `json:"model"`
	ManufacturingYear *int                `json:"manufacturing_year"`
	Price             *int64              `json:"price"`
	Type              *string             `json:"type"`
	FuelType          *string             `json:"fuel_type"`
	MileageType       *float64            `json:"mileage"`
	Attachments       []VehicleAttachment `json:"attachments"`
}

type VehicleAttachment struct {
	Id        *uuid.UUID `json:"id"`
	Path      *string    `json:"path"`
	MediaType *string    `json:"media_type"`
}

type VehicleComparisonResponse struct {
	Ids             []interface{}            `json:"ids"`
	ComparisonTable map[string][]interface{} `json:"comparison_table"`
}
