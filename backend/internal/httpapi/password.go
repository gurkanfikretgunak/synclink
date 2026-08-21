package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (s *Server) forgotPassword(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || strings.TrimSpace(in.Email) == "" {
		writeJSON(w, 400, map[string]string{"error": "invalid email"})
		return
	}
	token, found := s.auth.ForgotPassword(r.Context(), in.Email)
	// Always 200 so we do not leak whether the email exists.
	out := map[string]any{"ok": true}
	if found {
		out["resetToken"] = token
		out["expiresIn"] = 1800
	}
	writeJSON(w, 200, out)
}

func (s *Server) resetPassword(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email       string `json:"email"`
		Token       string `json:"token"`
		NewPassword string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	if err := s.auth.ResetPassword(r.Context(), in.Email, in.Token, in.NewPassword); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true})
}

func (s *Server) changePassword(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(r)
	if !ok {
		writeJSON(w, 401, map[string]string{"error": "not authenticated"})
		return
	}
	var in struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	if err := s.auth.ChangePassword(r.Context(), id, in.CurrentPassword, in.NewPassword); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true})
}
