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
)

func vehicleRouter(router *mux.Router) {
	subRouter := router.PathPrefix(subRouterVehicleV1Path).Subrouter()

	vehicleHandler := handlers.NewVehicleHandler()
	subRouter.Methods(http.MethodGet).Path(searchByModelNamePath).Handler(
		middleware.NilRequestResponseMw(
			vehicleHandler.GetVehiclesByModelName,
		),
	)
}
