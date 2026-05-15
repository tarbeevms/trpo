package app

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"taskflow/internal/models"
	"taskflow/internal/service"
	"taskflow/pkg/auth"
)

type Server struct {
	users    *service.UserService
	projects *service.ProjectService
	tasks    *service.TaskFacade
	reports  *service.ReportService
	auth     *auth.Manager
}

type pageData struct {
	Users          []models.User
	Projects       []models.Project
	Tasks          []models.Task
	ExportedReport string // Отчёт, сгенерированный экспортёром
	ExportFormat   string // Формат экспорта (html, json, xml)
	Error          string
	Message        string
	Authenticated  bool
	CurrentUser    models.User
}

func NewServer(users *service.UserService, projects *service.ProjectService, tasks *service.TaskFacade, reports *service.ReportService, authManager *auth.Manager) *Server {
	return &Server{users: users, projects: projects, tasks: tasks, reports: reports, auth: authManager}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.index)
	mux.HandleFunc("/me", s.me)
	mux.HandleFunc("/register", s.register)
	mux.HandleFunc("/login", s.login)
	mux.HandleFunc("/logout", s.logout)
	mux.HandleFunc("/projects", s.createProject)
	mux.HandleFunc("/projects/delete", s.deleteProject)
	mux.HandleFunc("/tasks", s.createTask)
	mux.HandleFunc("/tasks/status", s.changeTaskStatus)
	mux.HandleFunc("/tasks/delete", s.deleteTask)
	mux.HandleFunc("/report", s.report)
	return mux
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if _, ok := s.currentUser(r); ok {
		http.Redirect(w, r, "/me", http.StatusSeeOther)
		return
	}
	s.render(w, r, pageData{
		Error:   r.URL.Query().Get("error"),
		Message: r.URL.Query().Get("message"),
	})
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user, ok := s.currentUser(r)
	if !ok {
		http.Redirect(w, r, "/?error="+urlValue("login required"), http.StatusSeeOther)
		return
	}
	s.render(w, r, pageData{
		Error:         r.URL.Query().Get("error"),
		Message:       r.URL.Query().Get("message"),
		Authenticated: true,
		CurrentUser:   user,
	})
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
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
	deadline, err := parseDeadline(r.FormValue("deadline"))
	if err != nil {
		redirectWithErrorTo(w, r, "/me", err)
		return
	}

	task := models.Task{
		ProjectID:   projectID,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Status:      models.StatusNew,
		Priority:    models.Priority(r.FormValue("priority")),
		Deadline:    deadline,
		AssigneeID:  user.ID,
		AuditInfo:   models.AuditInfo{CreatedBy: user.ID},
		Tags:        parseTags(r.FormValue("tags")),
	}
	if _, err := s.tasks.Create(contextWithRequest(r), task); err != nil {
		redirectWithErrorTo(w, r, "/me", err)
		return
	}
	redirectWithMessageTo(w, r, "/me", "task created")
}

func (s *Server) changeTaskStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user, ok := s.currentUser(r)
	if !ok {
		redirectWithError(w, r, fmt.Errorf("login required"))
		return
	}
	taskID, err := parseInt64(r.FormValue("task_id"))
	if err != nil {
		redirectWithErrorTo(w, r, "/me", err)
		return
	}
	if err := s.tasks.ChangeStatusForAssignee(contextWithRequest(r), taskID, user.ID, models.Status(r.FormValue("status"))); err != nil {
		redirectWithErrorTo(w, r, "/me", err)
		return
	}
	redirectWithMessageTo(w, r, "/me", "task status changed")
}

func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user, ok := s.currentUser(r)
	if !ok {
		redirectWithError(w, r, fmt.Errorf("login required"))
		return
	}
	taskID, err := parseInt64(r.FormValue("task_id"))
	if err != nil {
		redirectWithErrorTo(w, r, "/me", err)
		return
	}
	if err := s.tasks.Delete(contextWithRequest(r), taskID, user.ID); err != nil {
		redirectWithErrorTo(w, r, "/me", err)
		return
	}
	redirectWithMessageTo(w, r, "/me", "task deleted")
}

func (s *Server) render(w http.ResponseWriter, r *http.Request, data pageData) {
	ctx := r.Context()
	if data.Authenticated {
		filter := models.TaskFilter{AssigneeID: data.CurrentUser.ID}
		projects, err := s.projects.ListByOwner(ctx, data.CurrentUser.ID)
		if err != nil {
			data.Error = err.Error()
		}
		tasks, err := s.tasks.List(ctx, filter)
		if err != nil {
			data.Error = err.Error()
		}
		// Генерируем HTML-отчёт только если он ещё не передан (например, из /report)
		if data.ExportedReport == "" {
			exported, err := s.reports.Build(ctx, filter, "html")
			if err != nil {
				data.Error = err.Error()
			}
			data.ExportedReport = string(exported)
			data.ExportFormat = "html"
		}
		data.Users = []models.User{data.CurrentUser}
		data.Projects = projects
		data.Tasks = tasks
	}

	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"date": func(value time.Time) string {
			if value.IsZero() {
				return ""
			}
			return value.Format("2006-01-02")
		},
		"safeHTML": func(html string) template.HTML {
			return template.HTML(html)
		},
	}).ParseFiles("frontend/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func redirectWithError(w http.ResponseWriter, r *http.Request, err error) {
	redirectWithErrorTo(w, r, "/", err)
}

func redirectWithMessage(w http.ResponseWriter, r *http.Request, message string) {
	redirectWithMessageTo(w, r, "/", message)
}

func redirectWithErrorTo(w http.ResponseWriter, r *http.Request, path string, err error) {
	http.Redirect(w, r, path+"?error="+urlValue(err.Error()), http.StatusSeeOther)
}

func redirectWithMessageTo(w http.ResponseWriter, r *http.Request, path string, message string) {
	http.Redirect(w, r, path+"?message="+urlValue(message), http.StatusSeeOther)
}

func urlValue(value string) string {
	return url.QueryEscape(value)
}

func parseInt64(value string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(value), 10, 64)
}

func parseDeadline(value string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(value))
}

func parseTags(value string) []models.Tag {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	tags := make([]models.Tag, 0, len(parts))
	for _, part := range parts {
		name := strings.TrimSpace(part)
		if name != "" {
			tags = append(tags, models.Tag{Name: name})
		}
	}
	return tags
}

func contextWithRequest(r *http.Request) context.Context {
	return r.Context()
}

func (s *Server) currentUser(r *http.Request) (models.User, bool) {
	cookie, err := r.Cookie("taskflow_token")
	if err != nil {
		return models.User{}, false
	}
	claims, err := s.auth.Validate(cookie.Value)
	if err != nil {
		return models.User{}, false
	}
	user, err := s.users.FindByID(r.Context(), claims.UserID)
	if err != nil {
		return models.User{}, false
	}
	return user, true
}
