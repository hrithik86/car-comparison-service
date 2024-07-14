package controllers

import "car-comparison-service/db/model"

type SuggestionsFactory struct{}

func (sf *SuggestionsFactory) GetSuggestionsController(vehicleType *model.VehicleType) ISuggestionsController {
	switch *vehicleType {
	case model.CAR:
		return NewCarSuggestionsController()
	case model.TRUCK:
		return NewTruckSuggestionsController()
	default:
		return NewDefaultSuggestionsController()
	}
}
