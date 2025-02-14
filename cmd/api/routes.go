package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users/registerUser", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/loginUser", app.userLoginHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/logoutUser", app.userLogoutHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/addMember", app.requirePermission("Admin", app.addMemberToGroupHandler))

	router.HandlerFunc(http.MethodPost, "/v1/buildings/addBuilding", app.requirePermission("Admin", app.addNewBuilding))
	router.HandlerFunc(http.MethodGet, "/v1/buildings/showBuildings", app.showBuildings)
	return app.recoverPanic(app.enableCORS(app.authenticate(router)))
}
