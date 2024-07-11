package routers

import (
	"car-comparison-service/service/api/middleware"
	"car-comparison-service/service/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	subRouterVehicleV1Path = "/api/v1/vehicle"
	searchByModelNamePath  = "/search"
	getVehicleByIdPath     = "/{id}"
	vehicleSuggestionsPath = "/{id}/suggestions"
)

func vehicleRouter(router *mux.Router) {
	subRouter := router.PathPrefix(subRouterVehicleV1Path).Subrouter()

	vehicleHandler := handlers.NewVehicleHandler()

	subRouter.Methods(http.MethodGet).Path(searchByModelNamePath).Handler(
		middleware.NilRequestResponseMw(
			vehicleHandler.GetVehiclesByModelName,
		),
	)

	subRouter.Methods(http.MethodGet).Path(getVehicleByIdPath).Handler(
		middleware.NilRequestResponseMw(
			vehicleHandler.GetVehicleById,
		),
	)

	subRouter.Methods(http.MethodGet).Path(vehicleSuggestionsPath).Handler(
		middleware.NilRequestResponseMw(
			vehicleHandler.GetVehicleSuggestions,
		),
	)
}
