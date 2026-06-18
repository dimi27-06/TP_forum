package dto

// ErrorResponse sert à renvoyer une erreur de façon claire :
// un code (genre 404) et un message lisible.
type ErrorResponse struct {
	Code    int
	Message string
}
