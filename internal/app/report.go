package app

import (
	"fmt"
	"net/http"

	"taskflow/internal/models"
)

func (s *Server) report(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user, ok := s.currentUser(r)
	if !ok {
		redirectWithError(w, r, fmt.Errorf("login required"))
		return
	}
	s.render(w, r, pageData{
		ReportType:     models.ReportType(r.URL.Query().Get("report_type")),
		FilterStatus:   r.URL.Query().Get("status"),
		FilterPriority: r.URL.Query().Get("priority"),
		Authenticated:  true,
		CurrentUser:    user,
	})
}
