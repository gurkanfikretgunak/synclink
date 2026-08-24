package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gurkanfikretgunak/synclink/backend/internal/auth"
	"github.com/gurkanfikretgunak/synclink/backend/internal/page"
)

func (s *Server) requireAdmin(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, ok := userID(r)
	if !ok || !s.auth.IsAdmin(id) {
		writeJSON(w, 403, map[string]string{"error": "forbidden"})
		return uuid.Nil, false
	}
	return id, true
}

func (s *Server) publicSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.auth.PublicSettings())
}

func (s *Server) adminMe(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	id, _ := userID(r)
	u, found := s.auth.User(id)
	if !found {
		writeJSON(w, 401, map[string]string{"error": "not authenticated"})
		return
	}
	writeJSON(w, 200, u.Info())
}

func (s *Server) adminUsers(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	writeJSON(w, 200, s.auth.ListUsers(r.Context()))
}

func (s *Server) adminPatchUser(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid user id"})
		return
	}
	var in struct {
		Role   *string `json:"role"`
		Status *string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	u, err := s.auth.UpdateUser(r.Context(), id, in.Role, in.Status)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, u)
}

func (s *Server) adminDeleteUser(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid user id"})
		return
	}
	if err := s.auth.DeleteUser(r.Context(), id); err != nil {
		writeErr(w, err)
		return
	}
	w.WriteHeader(204)
}

func (s *Server) adminSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	writeJSON(w, 200, s.auth.Settings())
}

func (s *Server) adminPutSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var in auth.Settings
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	writeJSON(w, 200, s.auth.UpdateSettings(in))
}

func (s *Server) adminPages(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	pages, err := s.pages.ListAll(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, pages)
}

func (s *Server) adminStats(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	users := s.auth.ListUsers(r.Context())
	pages, err := s.pages.ListAll(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	clicks, err := s.pages.SumClicks(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	pageClicks, err := s.pages.PageClicks(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	if pageClicks == nil {
		pageClicks = []page.PageClickStat{}
	}
	writeJSON(w, 200, map[string]any{"users": len(users), "pages": len(pages), "totalClicks": clicks, "pageClicks": pageClicks})
}

func (s *Server) adminPatchPage(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid page id"})
		return
	}
	var in struct {
		Verified *bool `json:"verified"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Verified == nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	p, err := s.pages.SetPageVerified(r.Context(), id, *in.Verified)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, p)
}
