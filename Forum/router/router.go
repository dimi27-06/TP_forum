// Le package "router" définit la liste des adresses (URL) du site et indique,
// pour chacune, quelle fonction du controller doit répondre.
// C'est un peu le standard téléphonique du forum : "telle adresse -> tel service".
package router

import (
	"forum/controllers"
	"forum/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterAssetRoutes sert les fichiers statiques (CSS, images, JavaScript...)
// présents dans le dossier "static". Toute URL commençant par /static/ y est dirigée.
func RegisterAssetRoutes(r *mux.Router) {
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))),
	)
}

// RegisterAuthRoutes branche les pages de compte (inscription, connexion, etc.).
// Note : "GET" = afficher la page, "POST" = envoyer un formulaire.
func RegisterAuthRoutes(r *mux.Router, c *controllers.AuthController) {
	r.HandleFunc("/register", c.ShowRegister).Methods("GET")
	r.HandleFunc("/register", c.Register).Methods("POST")
	r.HandleFunc("/login", c.ShowLogin).Methods("GET")
	r.HandleFunc("/login", c.Login).Methods("POST")
	r.HandleFunc("/logout", c.Logout).Methods("POST", "GET")
	// La page "/me" est protégée : il faut être connecté pour y accéder.
	r.Handle("/me", middleware.RequireAuth(http.HandlerFunc(c.Me))).Methods("GET")
}

// RegisterThreadRoutes branche les pages des fils de discussion.
func RegisterThreadRoutes(r *mux.Router, c *controllers.ThreadController) {
	// Pages publiques : tout le monde peut voir la liste et le détail des fils.
	r.HandleFunc("/", c.List).Methods("GET")
	r.HandleFunc("/threads", c.List).Methods("GET")
	// "{id:[0-9]+}" veut dire : un nombre dans l'URL (ex: /threads/42).
	r.HandleFunc("/threads/{id:[0-9]+}", c.Show).Methods("GET")

	// À partir d'ici, il faut être connecté : on enveloppe chaque page dans RequireAuth.
	auth := middleware.RequireAuth
	r.Handle("/threads/new", auth(http.HandlerFunc(c.ShowCreate))).Methods("GET")
	r.Handle("/threads/new", auth(http.HandlerFunc(c.Create))).Methods("POST")
	r.Handle("/threads/{id:[0-9]+}/edit", auth(http.HandlerFunc(c.ShowEdit))).Methods("GET")
	r.Handle("/threads/{id:[0-9]+}/edit", auth(http.HandlerFunc(c.Update))).Methods("POST")
	r.Handle("/threads/{id:[0-9]+}/delete", auth(http.HandlerFunc(c.Delete))).Methods("POST")
}

// RegisterMessageRoutes branche les actions sur les messages et les réactions.
func RegisterMessageRoutes(r *mux.Router, mc *controllers.MessageController, rc *controllers.ReactionController) {
	auth := middleware.RequireAuth

	// Toutes ces actions demandent d'etre connecté.
	r.Handle("/messages", auth(http.HandlerFunc(mc.Create))).Methods("POST")
	r.Handle("/messages/{id:[0-9]+}/edit", auth(http.HandlerFunc(mc.ShowEdit))).Methods("GET")
	r.Handle("/messages/{id:[0-9]+}/edit", auth(http.HandlerFunc(mc.Update))).Methods("POST")
	r.Handle("/messages/{id:[0-9]+}/delete", auth(http.HandlerFunc(mc.Delete))).Methods("POST")

	// Reactions (JSON endpoint)
	// "{type}" sera "like" ou "dislike" selon le bouton cliqué.
	r.Handle("/messages/{id:[0-9]+}/react/{type}", auth(http.HandlerFunc(rc.React))).Methods("POST")
}

// RegisterAdminRoutes branche le tableau de bord et les actions réservées aux admins.
// On utilise RequireAdmin : si on n'est pas admin, l'accès est refusé.
func RegisterAdminRoutes(r *mux.Router, c *controllers.AdminController) {
	admin := middleware.RequireAdmin

	r.Handle("/admin", admin(http.HandlerFunc(c.Dashboard))).Methods("GET")
	r.Handle("/admin/users/{id:[0-9]+}/ban", admin(http.HandlerFunc(c.BanUser))).Methods("POST")
	r.Handle("/admin/users/{id:[0-9]+}/unban", admin(http.HandlerFunc(c.UnbanUser))).Methods("POST")
	r.Handle("/admin/threads/{id:[0-9]+}/status", admin(http.HandlerFunc(c.UpdateThreadStatus))).Methods("POST")
	r.Handle("/admin/threads/{id:[0-9]+}/delete", admin(http.HandlerFunc(c.DeleteThread))).Methods("POST")
	r.Handle("/admin/messages/{id:[0-9]+}/delete", admin(http.HandlerFunc(c.DeleteMessage))).Methods("POST")
}
