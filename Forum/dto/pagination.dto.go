package dto

// PaginationRequest = les options choisies par l'utilisateur pour afficher une liste :
// quelle page, combien d'éléments, comment trier, et les filtres.
type PaginationRequest struct {
	Page   int
	Limit  int    // 10, 20, 30, ou 0 = tout
	Sort   string // "recent", "oldest", "popular"
	Tag    string // filtrer par tag
	Search string // filtrer par mot-clé
}

// PaginationMeta = les infos calculées pour afficher la barre de pagination
// (page courante, nombre total, y a-t-il une page avant/après, etc.).
type PaginationMeta struct {
	Page       int
	Limit      int
	Total      int  // nombre total d'éléments
	TotalPages int  // nombre total de pages
	HasPrev    bool // existe-t-il une page précédente ?
	HasNext    bool // existe-t-il une page suivante ?
}
