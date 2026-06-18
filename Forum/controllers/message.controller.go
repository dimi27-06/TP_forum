package controllers

import (
	"forum/dto"
	"forum/middleware"
	"forum/services"
	"forum/templates"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// MessageController gère la création, la modification et la suppression des messages.
type MessageController struct {
	messageService *services.MessageService
	tmpl           *templates.Manager
}

func InitMessageController(ms *services.MessageService, tmpl *templates.Manager) *MessageController {
	return &MessageController{messageService: ms, tmpl: tmpl}
}

// Create poste un nouveau message dans un fil de discussion.
func (c *MessageController) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	r.ParseForm()

	// L'id du fil arrive sous forme de texte : on le convertit en nombre.
	threadID, _ := strconv.Atoi(r.FormValue("thread_id"))
	req := dto.CreateMessageRequest{
		Content:  r.FormValue("content"),
		ThreadID: threadID,
	}

	_, err := c.messageService.Create(req, claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// On revient sur le fil pour voir le message ajouté.
	http.Redirect(w, r, "/threads/"+strconv.Itoa(threadID), http.StatusFound)
}

// ShowEdit affiche le formulaire de modification d'un message.
func (c *MessageController) ShowEdit(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	msgID, _ := strconv.Atoi(mux.Vars(r)["id"])

	msg, err := c.messageService.GetByID(msgID)
	if err != nil {
		http.Error(w, "Message introuvable", http.StatusNotFound)
		return
	}
	// Seul l'auteur du message ou un admin peut le modifier.
	if claims.Role != "admin" && msg.UserID != claims.UserID {
		http.Error(w, "Non autorisé", http.StatusForbidden)
		return
	}

	c.tmpl.Render(w, r, "thread/edit_message.html", map[string]interface{}{"Message": msg})
}

// Update enregistre la modification d'un message.
func (c *MessageController) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	msgID, _ := strconv.Atoi(mux.Vars(r)["id"])
	r.ParseForm()

	req := dto.UpdateMessageRequest{Content: r.FormValue("content")}
	// On récupère le message avant la modif pour connaitre son fil (pour la redirection).
	msg, _ := c.messageService.GetByID(msgID)

	if err := c.messageService.Update(msgID, claims.UserID, claims.Role, req); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// On retourne sur le fil contenant le message.
	threadID := 0
	if msg != nil {
		threadID = msg.ThreadID
	}
	http.Redirect(w, r, "/threads/"+strconv.Itoa(threadID), http.StatusFound)
}

// Delete supprime un message.
func (c *MessageController) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	msgID, _ := strconv.Atoi(mux.Vars(r)["id"])

	// On note d'abord le fil du message pour pouvoir y revenir après suppression.
	msg, _ := c.messageService.GetByID(msgID)
	threadID := 0
	if msg != nil {
		threadID = msg.ThreadID
	}

	if err := c.messageService.Delete(msgID, claims.UserID, claims.Role); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	http.Redirect(w, r, "/threads/"+strconv.Itoa(threadID), http.StatusFound)
}

// GetByID exposé pour les autres controllers
// Permet aux autres controllers (ex: admin) de réutiliser ce service.
func (c *MessageController) GetByIDService() *services.MessageService {
	return c.messageService
}
