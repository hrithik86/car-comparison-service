package api

import (
	"car-comparison-service/config"
	"car-comparison-service/service/routers"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	"github.com/urfave/negroni"
	"net/http"
	"strconv"
)

func StartAPI(c *cli.Context) error {
	router := routers.Router()
	startServer(router)
	return nil
}

func startServer(router *mux.Router) {
	server := negroni.New(negroni.NewRecovery())
	handlerFunc := router.ServeHTTP
	server.UseHandlerFunc(handlerFunc)
	portInfo := ":" + strconv.Itoa(int(config.Port()))
	queryTimeoutMiddleware := http.TimeoutHandler(server, config.QueryTimeout(), "Sorry, the request has timed out!")

	httpServer := &http.Server{
		Addr:    portInfo,
		Handler: queryTimeoutMiddleware,
	}
	fmt.Println("API server started at port " + portInfo)

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

}
