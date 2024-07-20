package routers

import (
	"car-comparison-service/service/api/middleware"
	"car-comparison-service/service/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	subRouterVehicleV1Path = "/api/v1/vehicle"

	searchByModelNamePath     = "/search"
	compareVehiclesPath       = "/compare"
	createVehiclePath         = "/create"
	getVehicleByIdPath        = "/{id}"
	vehicleSuggestionsPath    = "/{id}/suggestions"
	addVehicleAttachmentsPath = "/{id}/add-attachments"
	addVehicleFeaturesPath    = "/{id}/add-features"
)

func vehicleRouter(router *mux.Router) {
	subRouter := router.PathPrefix(subRouterVehicleV1Path).Subrouter()

	vehicleHandler := handlers.NewVehicleHandler()

	subRouter.Methods(http.MethodGet).Path(searchByModelNamePath).Handler(
		middleware.NilRequestResponseMw(
			vehicleHandler.GetVehiclesByModelName,
		),
	)

	subRouter.Methods(http.MethodGet).Path(compareVehiclesPath).Handler(
		middleware.RequestResponseMw(
			vehicleHandler.GetVehicleComparison,
		),
	)

	subRouter.Methods(http.MethodGet).Path(vehicleSuggestionsPath).Handler(
		middleware.NilRequestResponseMw(
			vehicleHandler.GetVehicleSuggestions,
		),
	)

	subRouter.Methods(http.MethodGet).Path(getVehicleByIdPath).Handler(
		middleware.NilRequestResponseMw(
			vehicleHandler.GetVehicleInfoById,
		),
	)

	subRouter.Methods(http.MethodPost).Path(createVehiclePath).Handler(
		middleware.RequestResponseMw(
			vehicleHandler.CreateVehicle,
		),
	)

	subRouter.Methods(http.MethodPost).Path(addVehicleAttachmentsPath).Handler(
		middleware.RequestResponseMw(
			vehicleHandler.AddVehicleAttachments,
		),
	)

	subRouter.Methods(http.MethodPost).Path(addVehicleFeaturesPath).Handler(
		middleware.RequestResponseMw(
			vehicleHandler.AddVehicleFeatures,
		),
	)
}
