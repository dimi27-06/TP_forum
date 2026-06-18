// Le package "templates" s'occupe d'afficher les pages HTML à partir de fichiers modèles.
// Un "template" est une page HTML à trous : on y injecte les données (titre, messages...).
package templates

import (
	"forum/middleware"
	"html/template" // le moteur de templates fourni par Go
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Manager garde une liste de petites fonctions utilisables directement dans les pages HTML.
type Manager struct {
	funcMap template.FuncMap
}

// NewManager crée le gestionnaire et définit les fonctions utilitaires pour les templates.
func NewManager() *Manager {
	return &Manager{
		funcMap: template.FuncMap{
			"add": func(a, b int) int { return a + b }, // addition (ex: numéro de page suivante)
			"sub": func(a, b int) int { return a - b }, // soustraction (ex: page précédente)
			// seq fabrique une liste 1, 2, 3, ... n. Pratique pour afficher les numéros de pages.
			"seq": func(n int) []int {
				s := make([]int, n)
				for i := range s {
					s[i] = i + 1
				}
				return s
			},
			// substr renvoie un morceau de texte (extrait), en gérant les bornes pour ne pas planter.
			"substr": func(s string, start, end int) string {
				runes := []rune(s) // on passe par les "runes" pour bien gérer les accents
				if start < 0 {
					start = 0
				}
				if end > len(runes) {
					end = len(runes)
				}
				if start >= end {
					return ""
				}
				return string(runes[start:end])
			},
		},
	}
}

// Render assemble la mise en page commune (layout) avec une page précise, puis l'envoie au navigateur.
func (m *Manager) Render(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	// On récupère l'utilisateur connecté pour le rendre disponible dans toutes les pages
	// (ça permet d'afficher son pseudo, ou le bouton "Admin" si besoin).
	claims := middleware.GetClaims(r)

	templateData := map[string]interface{}{
		"Claims": claims,
	}

	// On ajoute les données spécifiques à la page (la liste des fils, un message d'erreur, etc.).
	if d, ok := data.(map[string]interface{}); ok {
		for k, v := range d {
			templateData[k] = v
		}
	}

	// On cherche le dossier des templates (selon l'endroit d'où le programme est lancé).
	basePath := "./templates"
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		basePath = "templates"
	}

	// Chaque page utilise layout.html (l'habillage commun : entête, menu, pied de page)
	// auquel on greffe le contenu de la page demandée.
	layout := filepath.Join(basePath, "layout.html")
	page := filepath.Join(basePath, name)

	// On charge et prépare les deux fichiers ensemble.
	tmpl, err := template.New("layout.html").Funcs(m.funcMap).ParseFiles(layout, page)
	if err != nil {
		// Si un template est mal écrit ou introuvable, on le note dans les logs et on affiche une erreur.
		log.Printf("Template error (%s): %v", name, err)
		http.Error(w, "Erreur de rendu de la page", http.StatusInternalServerError)
		return
	}

	// On précise que c'est du HTML en UTF-8 (pour que les accents s'affichent bien), puis on génère la page.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, templateData); err != nil {
		log.Printf("Template execute error (%s): %v", name, err)
	}
}
