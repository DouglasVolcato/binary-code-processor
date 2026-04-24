package web

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed index.html
var content embed.FS

type Server struct {
	tmpl             *template.Template
	websocketURL     string
}

type pageData struct {
	WebSocketURL string
}

func NewServer(websocketURL string) (*Server, error) {
	tmpl, err := template.ParseFS(content, "index.html")
	if err != nil {
		return nil, err
	}

	return &Server{
		tmpl:         tmpl,
		websocketURL: websocketURL,
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.Execute(w, pageData{WebSocketURL: s.websocketURL}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
