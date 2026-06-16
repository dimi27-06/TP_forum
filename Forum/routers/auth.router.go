package routers

import (
	"exemple_api/controllers"
	"exemple_api/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// AuthProductRoutes enregistre les routes liees a l'authentification.
func AuthProductRoutes(r *mux.Router, authController *controllers.AuthControllers) {
	// Route publique : elle permet de se connecter et de recuperer un token JWT.
	r.HandleFunc("/login", authController.Login).Methods("POST")
	r.HandleFunc("/register", authController.Register).Methods("POST")

	// Route protegee : elle passe par le middleware pour verifier le token JWT.
	r.Handle("/me", middleware.AuthMiddleware(http.HandlerFunc(authController.Me))).Methods("GET")
}
