package router

import (
	"forum/controllers"
	"forum/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterAssetRoutes(r *mux.Router) {
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))),
	)
}

func RegisterAuthRoutes(r *mux.Router, c *controllers.AuthController) {
	r.HandleFunc("/register", c.ShowRegister).Methods("GET")
	r.HandleFunc("/register", c.Register).Methods("POST")
	r.HandleFunc("/login", c.ShowLogin).Methods("GET")
	r.HandleFunc("/login", c.Login).Methods("POST")
	r.HandleFunc("/logout", c.Logout).Methods("POST", "GET")
	r.Handle("/me", middleware.RequireAuth(http.HandlerFunc(c.Me))).Methods("GET")
}

func RegisterThreadRoutes(r *mux.Router, c *controllers.ThreadController) {
	r.HandleFunc("/", c.List).Methods("GET")
	r.HandleFunc("/threads", c.List).Methods("GET")
	r.HandleFunc("/threads/{id:[0-9]+}", c.Show).Methods("GET")

	auth := middleware.RequireAuth
	r.Handle("/threads/new", auth(http.HandlerFunc(c.ShowCreate))).Methods("GET")
	r.Handle("/threads/new", auth(http.HandlerFunc(c.Create))).Methods("POST")
	r.Handle("/threads/{id:[0-9]+}/edit", auth(http.HandlerFunc(c.ShowEdit))).Methods("GET")
	r.Handle("/threads/{id:[0-9]+}/edit", auth(http.HandlerFunc(c.Update))).Methods("POST")
	r.Handle("/threads/{id:[0-9]+}/delete", auth(http.HandlerFunc(c.Delete))).Methods("POST")
}

func RegisterMessageRoutes(r *mux.Router, mc *controllers.MessageController, rc *controllers.ReactionController) {
	auth := middleware.RequireAuth

	r.Handle("/messages", auth(http.HandlerFunc(mc.Create))).Methods("POST")
	r.Handle("/messages/{id:[0-9]+}/edit", auth(http.HandlerFunc(mc.ShowEdit))).Methods("GET")
	r.Handle("/messages/{id:[0-9]+}/edit", auth(http.HandlerFunc(mc.Update))).Methods("POST")
	r.Handle("/messages/{id:[0-9]+}/delete", auth(http.HandlerFunc(mc.Delete))).Methods("POST")
	r.Handle("/messages/{id:[0-9]+}/react/{type}", auth(http.HandlerFunc(rc.React))).Methods("POST")
}

func RegisterAdminRoutes(r *mux.Router, c *controllers.AdminController) {
	admin := middleware.RequireAdmin

	r.Handle("/admin", admin(http.HandlerFunc(c.Dashboard))).Methods("GET")
	r.Handle("/admin/users/{id:[0-9]+}/ban", admin(http.HandlerFunc(c.BanUser))).Methods("POST")
	r.Handle("/admin/users/{id:[0-9]+}/unban", admin(http.HandlerFunc(c.UnbanUser))).Methods("POST")
	r.Handle("/admin/threads/{id:[0-9]+}/status", admin(http.HandlerFunc(c.UpdateThreadStatus))).Methods("POST")
	r.Handle("/admin/threads/{id:[0-9]+}/delete", admin(http.HandlerFunc(c.DeleteThread))).Methods("POST")
	r.Handle("/admin/messages/{id:[0-9]+}/delete", admin(http.HandlerFunc(c.DeleteMessage))).Methods("POST")
}
