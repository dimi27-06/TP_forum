package models

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"nom"`
	Description string  `json:"description"`
	Price       float32 `json:"prix"`
	CategorieId int     `json:"categorie_id"`
	CreateAt    string  `json:"date_ajout"`
}