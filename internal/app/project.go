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
	}
	if _, err := s.projects.Create(contextWithRequest(r), project); err != nil {
		redirectWithErrorTo(w, r, "/me", err)
		return
	}
	redirectWithMessageTo(w, r, "/me", "project created")
}

func (s *Server) deleteProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user, ok := s.currentUser(r)
	if !ok {
		redirectWithError(w, r, fmt.Errorf("login required"))
		return
	}
	projectID, err := parseInt64(r.FormValue("project_id"))
	if err != nil {
		redirectWithErrorTo(w, r, "/me", err)
		return
	}
	if err := s.projects.Delete(contextWithRequest(r), projectID, user.ID); err != nil {
		redirectWithErrorTo(w, r, "/me", err)
		return
	}
	redirectWithMessageTo(w, r, "/me", "project deleted")
}
