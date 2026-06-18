// Le package "controllers" fait le lien entre le navigateur et les services.
// Il lit les formulaires, appelle le bon service, puis affiche une page ou redirige.
package controllers

import (
	"forum/dto"
	"forum/middleware"
	"forum/services"
	"forum/templates"
	"net/http"
	"time"
)

// AuthController gère l'inscription, la connexion, la déconnexion et le profil.
type AuthController struct {
	authService *services.AuthService
	tmpl        *templates.Manager
}

func InitAuthController(authService *services.AuthService, tmpl *templates.Manager) *AuthController {
	return &AuthController{authService: authService, tmpl: tmpl}
}

// ShowRegister affiche simplement la page du formulaire d'inscription.
func (c *AuthController) ShowRegister(w http.ResponseWriter, r *http.Request) {
	c.tmpl.Render(w, r, "auth/register.html", nil)
}

// Register traite l'envoi du formulaire d'inscription.
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // on lit les champs envoyés par le formulaire
	// On range les valeurs du formulaire dans une structure bien organisée.
	req := dto.RegisterRequest{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Confirm:  r.FormValue("confirm"),
	}

	// On demande au service de créer le compte. S'il y a un souci (email déjà pris, etc.),
	// on réaffiche le formulaire avec le message d'erreur.
	if err := c.authService.Register(req); err != nil {
		c.tmpl.Render(w, r, "auth/register.html", map[string]interface{}{
			"Error": err.Error(),
			"Form":  req, // on renvoie les valeurs pour ne pas tout retaper
		})
		return
	}

	// Inscription réussie : on envoie vers la page de connexion.
	http.Redirect(w, r, "/login?registered=1", http.StatusFound)
}

// ShowLogin affiche la page de connexion.
func (c *AuthController) ShowLogin(w http.ResponseWriter, r *http.Request) {
	c.tmpl.Render(w, r, "auth/login.html", map[string]interface{}{
		// "Registered" vaut vrai si on arrive juste après une inscription réussie.
		"Registered": r.URL.Query().Get("registered") == "1",
	})
}

// Login traite l'envoi du formulaire de connexion.
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	req := dto.LoginRequest{
		Identifier: r.FormValue("identifier"),
		Password:   r.FormValue("password"),
	}

	// Le service vérifie les identifiants et nous renvoie un jeton si tout est bon.
	_, token, err := c.authService.Login(req)
	if err != nil {
		c.tmpl.Render(w, r, "auth/login.html", map[string]interface{}{
			"Error":      err.Error(),
			"Identifier": req.Identifier,
		})
		return
	}

	// On dépose le jeton dans un cookie côté navigateur.
	// HttpOnly empêche le JavaScript d'y accéder (sécurité contre le vol de jeton).
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour), // le cookie dure 24h
	})

	// On redirige vers la page demandée au départ, sinon vers l'accueil.
	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/"
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

// Logout déconnecte l'utilisateur en effaçant son cookie.
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Astuce : on remet un cookie vide avec une date déjà passée,
	// ce qui force le navigateur à le supprimer.
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt_token",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

// Me affiche la page de profil de l'utilisateur connecté.
func (c *AuthController) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		// Pas connecté : direction la page de connexion.
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	c.tmpl.Render(w, r, "auth/profile.html", map[string]interface{}{
		"Claims": claims,
	})
}
