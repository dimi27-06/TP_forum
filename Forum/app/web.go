package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func registerWebRoutes(router *mux.Router) {
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("web"))))
	router.HandleFunc("/", servePage("web/index.html")).Methods(http.MethodGet)
	router.HandleFunc("/login", servePage("web/login.html")).Methods(http.MethodGet)
	router.HandleFunc("/register", servePage("web/register.html")).Methods(http.MethodGet)
	router.HandleFunc("/forum", servePage("web/forum.html")).Methods(http.MethodGet)
	router.HandleFunc("/forum/topic/{id}", servePage("web/topic.html")).Methods(http.MethodGet)
	router.HandleFunc("/guest", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/forum?mode=guest", http.StatusFound)
	}).Methods(http.MethodGet)
}

func servePage(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}
