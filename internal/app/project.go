package app

import (
	"fmt"
	"net/http"

	"taskflow/internal/models"
)

func (s *Server) createProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user, ok := s.currentUser(r)
	if !ok {
		redirectWithError(w, r, fmt.Errorf("login required"))
		return
	}
	project := models.Project{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		OwnerID:     user.ID,
		AuditInfo:   models.AuditInfo{CreatedBy: user.ID},
	}
	if _, err := s.projects.Create(contextWithRequest(r), project); err != nil {
		redirectWithErrorTo(w, r, "/me", err)
		return
	}
	redirectWithMessageTo(w, r, "/me", "project created")
}
