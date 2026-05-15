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
		// Для AJAX возвращаем JSON с ошибкой
		if isAjaxRequest(r) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error":"login required"}`))
			return
		}
		redirectWithError(w, r, fmt.Errorf("login required"))
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "html"
	}

	filter := models.TaskFilter{AssigneeID: user.ID}

	exported, err := s.reports.Build(r.Context(), filter, format)
	if err != nil {
		if isAjaxRequest(r) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error":"` + err.Error() + `"}`))
			return
		}
		redirectWithError(w, r, err)
		return
	}

	// Если AJAX-запрос — возвращаем только данные
	if isAjaxRequest(r) {
		if format == "html" {
			w.Header().Set("Content-Type", "text/html")
		} else if format == "json" {
			w.Header().Set("Content-Type", "application/json")
		} else if format == "xml" {
			w.Header().Set("Content-Type", "application/xml")
		}
		w.Write(exported)
		return
	}

	// Обычный запрос — рендерим страницу
	s.render(w, r, pageData{
		ExportedReport: string(exported),
		ExportFormat:   format,
		Authenticated:  true,
		CurrentUser:    user,
	})
}

// isAjaxRequest проверяет, является ли запрос AJAX (XMLHttpRequest)
func isAjaxRequest(r *http.Request) bool {
	return r.Header.Get("X-Requested-With") == "XMLHttpRequest" ||
		r.Header.Get("Accept") == "application/json"
}
