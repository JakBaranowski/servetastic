package servtatic

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	helpers "github.com/JakBaranowski/servetastic/helpers"
)

// NewHandler creates a new servtatic handler for serving http
func NewHandler(layoutDir string, contentDir string, layout string) Handler {
	log.Print("Creating new servtatic Handler")
	layoutFiles := helpers.GlobDir(layoutDir, "*.gohtml")
	contentFiles := helpers.GlobDir(contentDir, "*.gohtml")

	log.Print("Parsing templates")
	templates := parseLayoutAndContentTemplates(layoutFiles, contentFiles)
	if !baseTemplatesExist(templates, "home", "error") {
		panic("Need all base templates to work.")
	}
	templates["/"] = templates["home"]
	log.Print("Templates parsed")

	return Handler{
		templates: templates,
		layout:    layout,
	}
}

// Handler parses all layout files with each body file creating templates
// that are ready to be served to the user.
type Handler struct {
	templates map[string]*template.Template
	layout    string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page := filepath.Base(filepath.Clean(r.URL.Path))
	log.Printf("Page %v visited!", page)

	template := h.templates[page]
	if template == nil {
		log.Printf("Page not found")
		h.HandleError(w, r, 404)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	template.ExecuteTemplate(w, h.layout, nil)
}

// HandleError handles errors
func (h Handler) HandleError(w http.ResponseWriter, r *http.Request, errCode int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	h.templates["error"].ExecuteTemplate(w, h.layout, errCode)
}

func parseLayoutAndContentTemplates(layoutFiles []string, contentFiles []string) map[string]*template.Template {
	templates := make(map[string]*template.Template)
	for _, contentFile := range contentFiles {
		files := append(layoutFiles, contentFile)
		bodyFileName := helpers.StripPathAndExtension(contentFile)
		log.Printf("Parsing template %v", contentFile)
		template, err := template.ParseFiles(files...)
		if err != nil {
			panic("Failed parsing templates")
		}
		templates[bodyFileName] = template
	}
	return templates
}

func baseTemplatesExist(m map[string]*template.Template, keys ...string) bool {
	for _, key := range keys {
		if _, exist := m[key]; !exist {
			log.Printf("Base template \"%v\" does not exist", key)
			return exist
		}
	}
	return true
}
