package controllers

import (
	"forum/services"
	"forum/templates"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// AdminController regroupe les actions réservées aux administrateurs :
// bannir des membres, archiver/supprimer des sujets, supprimer des messages.
type AdminController struct {
	adminService   *services.AdminService
	threadService  *services.ThreadService
	messageService *services.MessageService
	tmpl           *templates.Manager
}

func InitAdminController(as *services.AdminService, ts *services.ThreadService, ms *services.MessageService, tmpl *templates.Manager) *AdminController {
	return &AdminController{adminService: as, threadService: ts, messageService: ms, tmpl: tmpl}
}

// Dashboard affiche le tableau de bord admin avec la liste des utilisateurs.
func (c *AdminController) Dashboard(w http.ResponseWriter, r *http.Request) {
	users, _ := c.adminService.GetAllUsers()
	c.tmpl.Render(w, r, "admin/dashboard.html", map[string]interface{}{
		"Users": users,
	})
}

// BanUser bannit un utilisateur (il ne pourra plus se connecter normalement).
func (c *AdminController) BanUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(mux.Vars(r)["id"])
	c.adminService.BanUser(userID)
	http.Redirect(w, r, "/admin", http.StatusFound)
}

// UnbanUser annule le bannissement d'un utilisateur.
func (c *AdminController) UnbanUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(mux.Vars(r)["id"])
	c.adminService.UnbanUser(userID)
	http.Redirect(w, r, "/admin", http.StatusFound)
}

// UpdateThreadStatus change le statut d'un fil (ouvrir, fermer, archiver).
func (c *AdminController) UpdateThreadStatus(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	threadID, _ := strconv.Atoi(mux.Vars(r)["id"])
	status := r.FormValue("status")
	c.threadService.UpdateStatus(threadID, status)
	http.Redirect(w, r, "/threads/"+strconv.Itoa(threadID), http.StatusFound)
}

// DeleteThread supprime un fil. Le "admin" passé en paramètre donne tous les droits.
func (c *AdminController) DeleteThread(w http.ResponseWriter, r *http.Request) {
	threadID, _ := strconv.Atoi(mux.Vars(r)["id"])
	c.threadService.Delete(threadID, 0, "admin")
	http.Redirect(w, r, "/", http.StatusFound)
}

// DeleteMessage supprime un message en tant qu'admin.
func (c *AdminController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	msgID, _ := strconv.Atoi(mux.Vars(r)["id"])
	// On récupère le fil du message avant suppression pour pouvoir y retourner.
	msg, _ := c.messageService.GetByID(msgID)
	threadID := 0
	if msg != nil {
		threadID = msg.ThreadID
	}
	c.messageService.Delete(msgID, 0, "admin")
	http.Redirect(w, r, "/threads/"+strconv.Itoa(threadID), http.StatusFound)
}
