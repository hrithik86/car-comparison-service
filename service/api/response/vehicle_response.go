package response

import "github.com/google/uuid"

type VehicleResponse struct {
	Id                *uuid.UUID          `json:"id"`
	Brand             *string             `json:"brand"`
	Model             *string             `json:"model"`
	ManufacturingYear *int                `json:"manufacturing_year"`
	Price             *int64              `json:"price"`
	Type              *string             `json:"type"`
	Attachments       []VehicleAttachment `json:"attachments"`
}

type VehicleAttachment struct {
	Id        *uuid.UUID `json:"id"`
	Path      *string    `json:"path"`
	MediaType *string    `json:"media_type"`
}