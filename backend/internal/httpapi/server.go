package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/gurkanfikretgunak/synclink/backend/internal/auth"
	"github.com/gurkanfikretgunak/synclink/backend/internal/page"
)

type ctxKey int

const userKey ctxKey = 1

type Server struct {
	auth  *auth.Service
	pages *page.Service
}

func New(authSvc *auth.Service, pages *page.Service) http.Handler {
	s := &Server{auth: authSvc, pages: pages}
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
	}))
	r.Get("/health/live", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "alive"})
	})
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", s.register)
		r.Post("/auth/login", s.login)
		r.Post("/auth/forgot-password", s.forgotPassword)
		r.Post("/auth/reset-password", s.resetPassword)
		r.Get("/public/pages/{slug}", s.publicPage)
		r.Group(func(r chi.Router) {
			r.Use(s.jwt)
			r.Get("/me", s.me)
			r.Put("/me/password", s.changePassword)
			r.Get("/me/page", s.getPage)
			r.Put("/me/page", s.upsertPage)
			r.Get("/me/page/links", s.listLinks)
			r.Post("/me/page/links", s.createLink)
			r.Put("/me/page/links/reorder", s.reorder)
			r.Patch("/me/page/links/{id}", s.updateLink)
			r.Delete("/me/page/links/{id}", s.deleteLink)
		})
	})
	return r
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, page.ErrNotFound):
		writeJSON(w, 404, map[string]any{"error": "Not Found", "message": err.Error(), "code": 404})
	case errors.Is(err, page.ErrConflict), errors.Is(err, auth.ErrExists):
		writeJSON(w, 409, map[string]any{"error": "Conflict", "message": err.Error(), "code": 409})
	case errors.Is(err, page.ErrValidation):
		writeJSON(w, 422, map[string]any{"error": "Unprocessable Entity", "message": err.Error(), "code": 422})
	case errors.Is(err, auth.ErrWeakPassword):
		writeJSON(w, 422, map[string]any{"error": "Unprocessable Entity", "message": err.Error(), "code": 422})
	case errors.Is(err, auth.ErrInvalidCreds), errors.Is(err, page.ErrUnauthorized):
		writeJSON(w, 401, map[string]any{"error": "Unauthorized", "message": err.Error(), "code": 401})
	default:
		writeJSON(w, 500, map[string]any{"error": "Internal Server Error", "message": "an internal error occurred", "code": 500})
	}
}

func (s *Server) jwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" {
			writeJSON(w, 401, map[string]string{"error": "missing authorization header"})
			return
		}
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			writeJSON(w, 401, map[string]string{"error": "invalid authorization format"})
			return
		}
		claims, err := s.auth.Parse(parts[1])
		if err != nil {
			writeJSON(w, 401, map[string]string{"error": "invalid token"})
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey, claims.UserID)))
	})
}

func userID(r *http.Request) (uuid.UUID, bool) {
	id, ok := r.Context().Value(userKey).(uuid.UUID)
	return id, ok
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" || len(in.Password) < 8 {
		writeJSON(w, 400, map[string]string{"error": "invalid email or password"})
		return
	}
	u, tok, err := s.auth.Register(r.Context(), strings.ToLower(strings.TrimSpace(in.Email)), in.Password)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, map[string]any{"token": tok, "user": map[string]any{"id": u.ID, "email": u.Email}})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	u, tok, err := s.auth.Login(r.Context(), strings.ToLower(strings.TrimSpace(in.Email)), in.Password)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, map[string]any{"token": tok, "user": map[string]any{"id": u.ID, "email": u.Email}})
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	id, ok := userID(r)
	if !ok {
		writeJSON(w, 401, map[string]string{"error": "not authenticated"})
		return
	}
	u, found := s.auth.User(id)
	if !found {
		writeJSON(w, 401, map[string]string{"error": "not authenticated"})
		return
	}
	writeJSON(w, 200, map[string]any{"id": u.ID, "email": u.Email})
}

func (s *Server) publicPage(w http.ResponseWriter, r *http.Request) {
	p, err := s.pages.GetPublicPage(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, p)
}

func (s *Server) getPage(w http.ResponseWriter, r *http.Request) {
	id, _ := userID(r)
	p, err := s.pages.GetMyPage(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, p)
}

func (s *Server) upsertPage(w http.ResponseWriter, r *http.Request) {
	id, _ := userID(r)
	var in page.UpsertPageInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	p, err := s.pages.UpsertPage(r.Context(), id, in)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, p)
}

func (s *Server) listLinks(w http.ResponseWriter, r *http.Request) {
	id, _ := userID(r)
	links, err := s.pages.ListLinks(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, links)
}

func (s *Server) createLink(w http.ResponseWriter, r *http.Request) {
	id, _ := userID(r)
	var in page.CreateLinkInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	link, err := s.pages.CreateLink(r.Context(), id, in)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, link)
}

func (s *Server) updateLink(w http.ResponseWriter, r *http.Request) {
	uid, _ := userID(r)
	lid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid link id"})
		return
	}
	var in page.UpdateLinkInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	link, err := s.pages.UpdateLink(r.Context(), uid, lid, in)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, link)
}

func (s *Server) deleteLink(w http.ResponseWriter, r *http.Request) {
	uid, _ := userID(r)
	lid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid link id"})
		return
	}
	if err := s.pages.DeleteLink(r.Context(), uid, lid); err != nil {
		writeErr(w, err)
		return
	}
	w.WriteHeader(204)
}

func (s *Server) reorder(w http.ResponseWriter, r *http.Request) {
	uid, _ := userID(r)
	var in struct {
		IDs []uuid.UUID `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	links, err := s.pages.ReorderLinks(r.Context(), uid, in.IDs)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, links)
}
