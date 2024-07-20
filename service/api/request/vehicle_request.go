package request

import (
	"car-comparison-service/db/model"
	"github.com/google/uuid"
)

type VehicleComparisonRequest struct {
	Ids                []uuid.UUID `json:"ids" validate:"max=3,dive,required"`
	HideCommonFeatures bool        `json:"hide_common_features"`
}

type CreateVehicleRequest struct {
	Model             *string            `json:"model"`
	Brand             *string            `json:"brand"`
	ManufacturingYear *int               `json:"manufacturing_year"`
	Type              *model.VehicleType `json:"vehicle_type"`
	Price             *int64             `json:"price"`
	FuelType          *model.FuelType    `json:"fuel_type"`
	Mileage           *float64           `json:"mileage"`
}

type BulkAddVehicleAttachmentsRequest struct {
	Name      *string          `json:"name"`
	Path      *string          `json:"path"`
	MediaType *model.MediaType `json:"media_type"`
}

type BulkAddVehicleFeaturesRequest struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}
