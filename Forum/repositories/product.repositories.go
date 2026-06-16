package repositories

import (
	"database/sql"
	"exemple_api/models"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type ProductRepository struct {
	db *sql.DB
}

func InitProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) CreateProduct(product models.Product) (int, error) {
	query := "INSERT INTO `produits`(`nom`, `description`, `prix`, `categorie_id`, `date_ajout`) VALUES (?,?,?,?,?);"

	sqlResult, sqlErr := r.db.Exec(query,
		product.Name,
		product.Description,
		product.Price,
		product.CategorieId,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	if sqlErr != nil {
		return -1, fmt.Errorf(" Erreur ajout produit - Erreur : \n\t %s", sqlErr.Error())
	}

	id, idErr := sqlResult.LastInsertId()
	if idErr != nil {
		return -1, fmt.Errorf(" Erreur ajout produit - Erreur récupération identifiant : \n\t %s", idErr.Error())
	}

	return int(id), nil
}

func (r *ProductRepository) ReadAll() ([]models.Product, error) {
	var listProducts []models.Product
	sqlResult, sqlErr := r.db.Query("SELECT * FROM `produits`;")
	if sqlErr != nil {
		return listProducts, fmt.Errorf(" Erreur récupération produit - Erreur : \n\t %s", sqlErr.Error())
	}

	defer sqlResult.Close()

	for sqlResult.Next() {
		var product models.Product
		errScan := sqlResult.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.CategorieId, &product.CreateAt)
		if errScan != nil {
			return nil, errScan
		}
		listProducts = append(listProducts, product)
	}

	return listProducts, nil
}

func (r *ProductRepository) ReadById(id int) (models.Product, error) {
	var product models.Product
	sqlErr := r.db.QueryRow("SELECT * FROM `produits` WHERE `produits`.id = ?;", id).
		Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.CategorieId, &product.CreateAt)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return models.Product{}, nil
		}
		return models.Product{}, fmt.Errorf(" Erreur récupération produit - Erreur : \n\t %s", sqlErr.Error())
	}

	return product, nil
}

func (r *ProductRepository) UpdateProductById(product models.Product) error {
	query := "UPDATE `produits` SET `nom`=?,`description`=?,`prix`=?,`categorie_id`=? WHERE `produits`.id=?;"

	sqlResult, sqlErr := r.db.Exec(query,
		product.Name,
		product.Description,
		product.Price,
		product.CategorieId,
		product.Id)

	if sqlErr != nil {
		return fmt.Errorf(" Erreur modification produit - Erreur : \n\t %s", sqlErr.Error())
	}

	if nbrRow, _ := sqlResult.RowsAffected(); nbrRow <= 0 {
		return fmt.Errorf(" Erreur modification produit - Aucune ligne modifiée")
	}

	return nil
}

func (r *ProductRepository) DeleteProductById(id int) error {
	sqlResult, sqlErr := r.db.Exec("DELETE FROM `produits` WHERE `produits`.id=?;", id)
	if sqlErr != nil {
		return fmt.Errorf(" Erreur suppression produit - Erreur : \n\t %s", sqlErr.Error())
	}

	if nbrRow, _ := sqlResult.RowsAffected(); nbrRow <= 0 {
		return fmt.Errorf(" Erreur suppression produit - Aucun produit supprime")
	}

	return nil
}
