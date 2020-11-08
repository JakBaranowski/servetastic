package servtatic

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// NewHandler creates a new servtatic handler for serving http
func NewHandler(layoutDir string, bodyDir string, layout string) Handler {
	log.Print("Creating new servtatic Handler")
	templates := make(map[string]*template.Template)
	layoutFiles := glob(layoutDir, "*.gohtml")
	bodyFiles := glob(bodyDir, "*.gohtml")

	log.Print("Parsing templates")
	for _, bodyFile := range bodyFiles {
		files := append(layoutFiles, bodyFile)
		bodyFileName := stripPathAndExtension(bodyFile)
		log.Printf("Parsing template %v", bodyFile)
		template, err := template.ParseFiles(files...)
		if err != nil {
			panic("Failed parsing templates")
		}
		templates[bodyFileName] = template
	}
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

	if page == "/" {
		page = "home"
	}

	template := h.templates[page]
	if template == nil {
		log.Printf("Page not found")
		return
	}

	template.ExecuteTemplate(w, h.layout, nil)
}

func glob(dir string, pattern string) []string {
	log.Printf("Globbing files at %v", dir)
	files, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		panic("Failed globing files")
	}
	log.Printf("Files \"%v\" globbed at %v", files, dir)
	return files
}

func stripPathAndExtension(s string) string {
	return strings.TrimSuffix(filepath.Base(s), filepath.Ext(s))
}
