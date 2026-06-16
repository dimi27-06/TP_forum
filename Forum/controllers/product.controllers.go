package controllers

import (
	"encoding/json"
	"exemple_api/helper"
	"exemple_api/models"
	"exemple_api/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProductControllers struct {
	service *services.ProductService
}

func InitProductController(service *services.ProductService) *ProductControllers {
	return &ProductControllers{service: service}
}

func readProductId(r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)["id"])
}

func (c *ProductControllers) Create(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "JSON invalide")
		return
	}

	productId, productErr := c.service.Create(newProduct)
	if productErr != nil {
		helper.WriteError(w, http.StatusBadRequest, productErr.Error())
		return
	}

	product, productErr := c.service.ReadById(productId)
	if productErr != nil {
		helper.WriteError(w, http.StatusInternalServerError, productErr.Error())
		return
	}

	helper.WriteJSON(w, http.StatusCreated, product)
}

func (c *ProductControllers) ReadAll(w http.ResponseWriter, r *http.Request) {
	productList, productErr := c.service.ReadAll()
	if productErr != nil {
		helper.WriteError(w, http.StatusInternalServerError, productErr.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, productList)
}

func (c *ProductControllers) ReadById(w http.ResponseWriter, r *http.Request) {
	idProduct, idProductErr := readProductId(r)
	if idProductErr != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant produit invalide")
		return
	}

	product, productErr := c.service.ReadById(idProduct)
	if productErr != nil {
		helper.WriteError(w, http.StatusInternalServerError, productErr.Error())
		return
	}
	if product.Id == 0 {
		helper.WriteError(w, http.StatusNotFound, "Produit introuvable")
		return
	}

	helper.WriteJSON(w, http.StatusOK, product)
}

func (c *ProductControllers) UpdateById(w http.ResponseWriter, r *http.Request) {
	idProduct, idProductErr := readProductId(r)
	if idProductErr != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant produit invalide")
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "JSON invalide")
		return
	}
	product.Id = idProduct

	productErr := c.service.UpdateById(product)
	if productErr != nil {
		helper.WriteError(w, http.StatusBadRequest, productErr.Error())
		return
	}

	updatedProduct, productErr := c.service.ReadById(idProduct)
	if productErr != nil {
		helper.WriteError(w, http.StatusInternalServerError, productErr.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, updatedProduct)
}

func (c *ProductControllers) DeleteById(w http.ResponseWriter, r *http.Request) {
	idProduct, idProductErr := readProductId(r)
	if idProductErr != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant produit invalide")
		return
	}

	productErr := c.service.DeleteById(idProduct)
	if productErr != nil {
		helper.WriteError(w, http.StatusBadRequest, productErr.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Produit supprime",
	})
}
