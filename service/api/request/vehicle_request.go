package request

import "github.com/google/uuid"

type VehicleComparisonRequest struct {
	Ids                []uuid.UUID `json:"ids" validate:"max=3,dive,required"`
	HideCommonFeatures bool        `json:"hide_common_features"`
}
