package services

import (
	"exemple_api/models"
	"exemple_api/repositories"
	"fmt"
)

type ProductService struct {
	// liste des repositories utilisés dans le service
	productRepository *repositories.ProductRepository
}

func InitProductService(productRepository *repositories.ProductRepository) *ProductService {
	return &ProductService{productRepository: productRepository}
}

func (s *ProductService) Create(product models.Product) (int, error) {
	if product.Name == "" || product.Description == "" || product.Price < 0 || product.CategorieId < 0 {
		return -1, fmt.Errorf(" Erreur ajout produit - Données manquantes ou invalides")
	}

	productId, prodproductErr := s.productRepository.CreateProduct(product)
	if prodproductErr != nil {
		return -1, prodproductErr
	}

	return productId, nil
}

func (s *ProductService) ReadAll() ([]models.Product, error) {
	productsList, productsErr := s.productRepository.ReadAll()
	if productsErr != nil {
		return nil, productsErr
	}

	return productsList, nil
}

func (s *ProductService) ReadById(idProduct int) (models.Product, error) {
	if idProduct <= 0 {
		return models.Product{}, fmt.Errorf(" Erreur récupération produit - identifiant invalide : %d", idProduct)
	}

	product, productErr := s.productRepository.ReadById(idProduct)
	if productErr != nil {
		return models.Product{}, productErr
	}

	return product, nil
}

func (s *ProductService) UpdateById(product models.Product) error {
	if product.Id <= 0 || product.Name == "" || product.Description == "" || product.Price < 0 || product.CategorieId < 0 {
		return fmt.Errorf(" Erreur modification produit - Donnees manquantes ou invalides")
	}

	return s.productRepository.UpdateProductById(product)
}

func (s *ProductService) DeleteById(idProduct int) error {
	if idProduct <= 0 {
		return fmt.Errorf(" Erreur suppression produit - identifiant invalide : %d", idProduct)
	}

	return s.productRepository.DeleteProductById(idProduct)
}
