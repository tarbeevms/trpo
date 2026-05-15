package app

import (
	"net/http"
	"time"

	"taskflow/internal/models"
)

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user, err := s.users.Register(contextWithRequest(r), r.FormValue("login"), r.FormValue("password"))
	if err != nil {
		redirectWithError(w, r, err)
		return
	}
	if err := s.setAuthCookie(w, user); err != nil {
		redirectWithError(w, r, err)
		return
	}
	redirectWithMessageTo(w, r, "/me", "registration completed")
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user, err := s.users.Login(contextWithRequest(r), r.FormValue("login"), r.FormValue("password"))
	if err != nil {
		redirectWithError(w, r, err)
		return
	}
	if err := s.setAuthCookie(w, user); err != nil {
		redirectWithError(w, r, err)
		return
	}
	redirectWithMessageTo(w, r, "/me", "login completed")
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "taskflow_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	redirectWithMessage(w, r, "logout completed")
}

func (s *Server) setAuthCookie(w http.ResponseWriter, user models.User) error {
	token, err := s.auth.Generate(user.ID, user.Login)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "taskflow_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}
