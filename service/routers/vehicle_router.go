package routers

import (
	"car-comparison-service/service/api/middleware"
	"car-comparison-service/service/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	subRouterVehicleV1Path = "/api/v1/vehicle"

	getVehicleByIdPath     = "/{id}"
	vehicleSuggestionsPath = "/{id}/suggestions"
	searchByModelNamePath  = "/search"
	compareVehiclesPath    = "/compare"
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
			vehicleHandler.GetVehicleById,
		),
	)
}
