package models

// Tag représente une étiquette (ex: "Chasse", "Pêche") servant à classer les sujets.
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
