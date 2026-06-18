package controllers

import (
	"forum/dto"
	"forum/middleware"
	"forum/services"
	"forum/templates"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type ThreadController struct {
	threadService  *services.ThreadService
	messageService *services.MessageService
	tmpl           *templates.Manager
}

func InitThreadController(ts *services.ThreadService, ms *services.MessageService, tmpl *templates.Manager) *ThreadController {
	return &ThreadController{threadService: ts, messageService: ms, tmpl: tmpl}
}

func (c *ThreadController) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit := 10
	if page < 1 {
		page = 1
	}
	if rawLimit := q.Get("limit"); rawLimit != "" {
		if parsedLimit, err := strconv.Atoi(rawLimit); err == nil {
			limit = parsedLimit
		}
	}

	req := dto.PaginationRequest{
		Page:   page,
		Limit:  limit,
		Sort:   q.Get("sort"),
		Tag:    q.Get("tag"),
		Search: q.Get("search"),
	}

	threads, meta, err := c.threadService.List(req)
	tags, _ := c.threadService.GetAllTags()
	if err != nil {
		c.tmpl.Render(w, r, "thread/list.html", map[string]interface{}{"Error": err.Error()})
		return
	}

	c.tmpl.Render(w, r, "thread/list.html", map[string]interface{}{
		"Threads": threads,
		"Meta":    meta,
		"Tags":    tags,
		"Query":   req,
	})
}

func (c *ThreadController) Show(w http.ResponseWriter, r *http.Request) {
	id := threadID(r)
	claims := middleware.GetClaims(r)

	thread, err := c.threadService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit := 10
	sort := q.Get("sort")
	if page < 1 {
		page = 1
	}
	if rawLimit := q.Get("limit"); rawLimit != "" {
		if parsedLimit, err := strconv.Atoi(rawLimit); err == nil {
			limit = parsedLimit
		}
	}

	currentUserID := 0
	if claims != nil {
		currentUserID = claims.UserID
	}

	messages, meta, _ := c.messageService.GetByThread(id, page, limit, sort, currentUserID)

	c.tmpl.Render(w, r, "thread/detail.html", map[string]interface{}{
		"Thread":   thread,
		"Messages": messages,
		"Meta":     meta,
		"Sort":     sort,
		"Limit":    limit,
	})
}

func (c *ThreadController) ShowCreate(w http.ResponseWriter, r *http.Request) {
	tags, _ := c.threadService.GetAllTags()
	c.tmpl.Render(w, r, "thread/create.html", map[string]interface{}{"Tags": tags})
}

func (c *ThreadController) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	r.ParseForm()

	req := dto.CreateThreadRequest{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
		Tags:    r.Form["tags"],
	}

	if free := r.FormValue("free_tags"); free != "" {
		for _, t := range strings.Split(free, ",") {
			req.Tags = append(req.Tags, strings.TrimSpace(t))
		}
	}

	id, err := c.threadService.Create(req, claims.UserID)
	if err != nil {
		tags, _ := c.threadService.GetAllTags()
		c.tmpl.Render(w, r, "thread/create.html", map[string]interface{}{
			"Error": err.Error(), "Form": req, "Tags": tags,
		})
		return
	}

	http.Redirect(w, r, "/threads/"+strconv.Itoa(id), http.StatusFound)
}

func (c *ThreadController) ShowEdit(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	id := threadID(r)

	thread, err := c.threadService.GetByIDAdmin(id)
	if err != nil {
		http.Error(w, "Fil introuvable", http.StatusNotFound)
		return
	}
	if claims.Role != "admin" && thread.UserID != claims.UserID {
		http.Error(w, "Non autorisé", http.StatusForbidden)
		return
	}

	tags, _ := c.threadService.GetAllTags()
	c.tmpl.Render(w, r, "thread/edit.html", map[string]interface{}{
		"Thread": thread, "Tags": tags,
	})
}

func (c *ThreadController) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	id := threadID(r)
	r.ParseForm()

	req := dto.UpdateThreadRequest{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
		Status:  r.FormValue("status"),
		Tags:    r.Form["tags"],
	}

	if err := c.threadService.Update(id, claims.UserID, claims.Role, req); err != nil {
		thread, _ := c.threadService.GetByIDAdmin(id)
		tags, _ := c.threadService.GetAllTags()
		c.tmpl.Render(w, r, "thread/edit.html", map[string]interface{}{
			"Error": err.Error(), "Thread": thread, "Tags": tags,
		})
		return
	}

	http.Redirect(w, r, "/threads/"+strconv.Itoa(id), http.StatusFound)
}

func (c *ThreadController) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	id := threadID(r)

	if err := c.threadService.Delete(id, claims.UserID, claims.Role); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func threadID(r *http.Request) int {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	return id
}
