package handlers

import (
	"car-comparison-service/errors"
	"car-comparison-service/serdes"
	"car-comparison-service/service/api/request"
	"car-comparison-service/service/controllers"
	"car-comparison-service/service/view"
	"context"
	"github.com/google/uuid"
	"net/http"
)

type VehicleHandler struct {
	vehicleController controllers.Vehicle
}

func NewVehicleHandler() *VehicleHandler {
	return &VehicleHandler{
		vehicleController: controllers.InitializeVehicleController(),
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
	return serdes.NewHttpResponse(http.StatusOK, view.CreateVehicleSearchResponse(response)), nil
}

func (v *VehicleHandler) GetVehicleInfoById(ctx context.Context, r serdes.Request[serdes.NilBody]) (serdes.Response, error) {
	vehicleId, err := uuid.Parse(r.Param("id"))
	if err != nil {
		return nil, errors.INVALID_UUID
	}

	response, err := v.vehicleController.GetVehicleInfoById(ctx, vehicleId)
	if err != nil {
		return nil, err
	}
	return serdes.NewHttpResponse(http.StatusOK, view.CreateVehicleFeaturesResponse(response)), nil
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

func (v *VehicleHandler) GetVehicleComparison(ctx context.Context, r serdes.Request[request.VehicleComparisonRequest]) (serdes.Response, error) {
	response, err := v.vehicleController.GetVehicleComparison(ctx, r.Body())
	if err != nil {
		return nil, err
	}
	return serdes.NewHttpResponse(http.StatusOK, view.CreateVehicleComparisonResponse(response)), nil
}

func (v *VehicleHandler) CreateVehicle(ctx context.Context, r serdes.Request[request.CreateVehicleRequest]) (serdes.Response, error) {
	response, err := v.vehicleController.CreateVehicle(ctx, r.Body())
	if err != nil {
		return nil, err
	}
	return serdes.NewHttpResponse(http.StatusOK, response), nil
}

func (v *VehicleHandler) AddVehicleAttachments(ctx context.Context, r serdes.Request[[]request.BulkAddVehicleAttachmentsRequest]) (serdes.Response, error) {
	vehicleId, err := uuid.Parse(r.Param("id"))
	if err != nil {
		return nil, errors.INVALID_UUID
	}

	response, err := v.vehicleController.AddVehicleAttachments(ctx, vehicleId, r.Body())
	if err != nil {
		return nil, err
	}
	return serdes.NewHttpResponse(http.StatusOK, response), nil
}

func (v *VehicleHandler) AddVehicleFeatures(ctx context.Context, r serdes.Request[[]request.BulkAddVehicleFeaturesRequest]) (serdes.Response, error) {
	vehicleId, err := uuid.Parse(r.Param("id"))
	if err != nil {
		return nil, errors.INVALID_UUID
	}

	response, err := v.vehicleController.AddVehicleFeatures(ctx, vehicleId, r.Body())
	if err != nil {
		return nil, err
	}
	return serdes.NewHttpResponse(http.StatusOK, response), nil
}
