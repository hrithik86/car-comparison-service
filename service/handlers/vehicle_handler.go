package handlers

import (
	"car-comparison-service/errors"
	"car-comparison-service/serdes"
	"car-comparison-service/service/controllers"
	"context"
	"github.com/google/uuid"
	"net/http"
)

type VehicleHandler struct {
	vehicleController *controllers.VehicleController
}

func NewVehicleHandler() *VehicleHandler {
	return &VehicleHandler{
		vehicleController: controllers.NewVehicleController(),
	}
}

func (v *VehicleHandler) GetVehiclesByModelName(ctx context.Context, r serdes.Request[serdes.NilBody]) (serdes.Response, error) {
	modelName := r.QueryParams().Get("modelName")
	if len(modelName) == 0 {
		return nil, errors.BAD_REQUEST
	}
	response, err := v.vehicleController.GetVehiclesByModelName(ctx, modelName)
	if err != nil {
		return nil, err
	}
	return serdes.NewHttpResponse(http.StatusOK, response), nil
}

func (v *VehicleHandler) GetVehicleById(ctx context.Context, r serdes.Request[serdes.NilBody]) (serdes.Response, error) {
	vehicleId, err := uuid.Parse(r.Param("id"))
	if err != nil {
		return nil, errors.INVALID_UUID
	}

	response, err := v.vehicleController.GetVehicleById(ctx, vehicleId)
	if err != nil {
		return nil, err
	}
	return serdes.NewHttpResponse(http.StatusOK, response), nil
}

func (v *VehicleHandler) GetVehicleSuggestions(ctx context.Context, r serdes.Request[serdes.NilBody]) (serdes.Response, error) {
	vehicleId, err := uuid.Parse(r.Param("id"))
	if err != nil {
		return nil, errors.INVALID_UUID
	}

	response, err := v.vehicleController.GetVehicleSuggestions(ctx, vehicleId)
	if err != nil {
		return nil, err
	}
	return serdes.NewHttpResponse(http.StatusOK, response), nil
}