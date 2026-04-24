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
	websocketPort    string
}

type pageData struct {
	WebSocketURL  string
	WebSocketPort string
}

func NewServer(websocketURL string, websocketPort string) (*Server, error) {
	tmpl, err := template.ParseFS(content, "index.html")
	if err != nil {
		return nil, err
	}

	return &Server{
		tmpl:          tmpl,
		websocketURL:  websocketURL,
		websocketPort: websocketPort,
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.Execute(w, pageData{WebSocketURL: s.websocketURL, WebSocketPort: s.websocketPort}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
