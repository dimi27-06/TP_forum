package routers

import (
	"exemple_api/controllers"

	"github.com/gorilla/mux"
)

func RegisterProductRoutes(r *mux.Router, productController *controllers.ProductControllers) {
	r.HandleFunc("/products", productController.ReadAll).Methods("GET")
	r.HandleFunc("/products/{id}", productController.ReadById).Methods("GET")
	r.HandleFunc("/products", productController.Create).Methods("POST")
	r.HandleFunc("/products/{id}", productController.UpdateById).Methods("PUT")
	r.HandleFunc("/products/{id}", productController.DeleteById).Methods("DELETE")
}
