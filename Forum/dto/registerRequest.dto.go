package dto

type RegisterRequestDto struct {
	Nom          string `json:"nom"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Bio          string `json:"bio"`
	Localisation string `json:"localisation"`
}
