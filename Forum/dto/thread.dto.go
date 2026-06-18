package dto

// CreateThreadRequest = les infos pour créer un nouveau fil de discussion.
type CreateThreadRequest struct {
	Title   string
	Content string
	Tags    []string // une liste de tags (les crochets [] veulent dire "plusieurs")
}

// UpdateThreadRequest = les infos pour modifier un fil existant.
type UpdateThreadRequest struct {
	Title   string
	Content string
	Status  string // open, closed ou archived
	Tags    []string
}
