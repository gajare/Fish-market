package router

import (
	"log"
	"net/http"

	"github.com/gajare/Fish-market/controller"
	"github.com/gajare/Fish-market/middleware"
	"github.com/gorilla/mux"
)

func New(c *controller.UserController) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	// global middlewares
	r.Use(middleware.RequestID) // injects req_id
	r.Use(middleware.AccessLog) // logs every request

	api := r.PathPrefix("/api/v1").Subrouter()

	// Public
	api.HandleFunc("/auth/login", c.Login).Methods(http.MethodPost)
	api.HandleFunc("/auth/register", c.Register).Methods(http.MethodPost)

	// Protected
	protected := api.NewRoute().Subrouter()
	protected.Use(middleware.Auth)

	// Admin-only
	admin := protected.NewRoute().Subrouter()
	admin.Use(middleware.AllowRoles("admin"))

	ac := controller.NewAdminController()
	admin.HandleFunc("/admin/log-level", ac.GetLogLevel).Methods(http.MethodGet)
	admin.HandleFunc("/admin/log-level", ac.SetLogLevel).Methods(http.MethodPut)

	admin.HandleFunc("/users", c.List).Methods(http.MethodGet)
	admin.HandleFunc("/users/{id}", c.Delete).Methods(http.MethodDelete)

	// Self or admin
	protected.HandleFunc("/users/{id}", c.GetByID).Methods(http.MethodGet)
	protected.HandleFunc("/users/{id}", c.Update).Methods(http.MethodPatch)

	// Log all mounted routes at startup
	_ = r.Walk(func(rt *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		path, _ := rt.GetPathTemplate()
		methods, _ := rt.GetMethods()
		if len(methods) == 0 {
			methods = []string{"ANY"}
		}
		if path == "" {
			path = "(no path template)"
		}
		log.Printf("route: %-6v %s", methods, path)
		return nil
	})

	return r
}

//ok
