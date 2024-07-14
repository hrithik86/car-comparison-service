package response

import "github.com/google/uuid"

type VehicleBaseResponse struct {
	Id                *uuid.UUID `json:"id"`
	Brand             *string    `json:"brand"`
	Model             *string    `json:"model"`
	ManufacturingYear *int       `json:"manufacturing_year"`
	Price             *int64     `json:"price"`
	Type              *string    `json:"type"`
	FuelType          *string    `json:"fuel_type"`
	MileageType       *float64   `json:"mileage"`
}

type VehicleWithAttachmentsResponse struct {
	VehicleBaseResponse
	Attachments []VehicleAttachment `json:"attachments"`
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

type VehicleWithFeaturesResponse struct {
	VehicleBaseResponse
	Features []VehicleFeatures `json:"features"`
}

type VehicleFeatures struct {
	FeatureId *uuid.UUID `json:"feature_id"`
	Key       *string    `json:"key"`
	Value     *string    `json:"value"`
}
