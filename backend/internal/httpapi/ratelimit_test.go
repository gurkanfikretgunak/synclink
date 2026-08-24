package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gurkanfikretgunak/synclink/backend/internal/auth"
	"github.com/gurkanfikretgunak/synclink/backend/internal/page"
)

func TestPublicClickRateLimitHeaders(t *testing.T) {
	authSvc := auth.NewService()
	pages := page.NewService(page.NewMemoryStore())
	h := New(authSvc, pages)

	reg := httptest.NewRecorder()
	h.ServeHTTP(reg, httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(`{"email":"rl@b.com","password":"password1"}`)))
	if reg.Code != 201 {
		t.Fatalf("register %d %s", reg.Code, reg.Body.String())
	}
	var login struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(reg.Body.Bytes(), &login); err != nil || login.Token == "" {
		t.Fatalf("token %#v", login)
	}
	put := httptest.NewRequest(http.MethodPut, "/api/v1/me/page", bytes.NewBufferString(`{"slug":"rl","displayName":"RL"}`))
	put.Header.Set("Authorization", "Bearer "+login.Token)
	put.Header.Set("Content-Type", "application/json")
	pw := httptest.NewRecorder()
	h.ServeHTTP(pw, put)
	if pw.Code != 200 {
		t.Fatalf("page %d %s", pw.Code, pw.Body.String())
	}
	cr := httptest.NewRequest(http.MethodPost, "/api/v1/me/page/links", bytes.NewBufferString(`{"title":"Site","url":"https://example.com"}`))
	cr.Header.Set("Authorization", "Bearer "+login.Token)
	cr.Header.Set("Content-Type", "application/json")
	cw := httptest.NewRecorder()
	h.ServeHTTP(cw, cr)
	if cw.Code != 201 {
		t.Fatalf("link %d %s", cw.Code, cw.Body.String())
	}
	var link page.LinkDTO
	if err := json.Unmarshal(cw.Body.Bytes(), &link); err != nil {
		t.Fatal(err)
	}

	click := func(ip string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/public/pages/rl/links/"+link.ID.String()+"/click", nil)
		req.Header.Set("X-Forwarded-For", ip)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		return w
	}

	w1 := click("203.0.113.10")
	if w1.Code != 200 {
		t.Fatalf("click1 %d %s", w1.Code, w1.Body.String())
	}
	if w1.Header().Get("X-RateLimit-Limit") != "60" {
		t.Fatalf("limit header %q", w1.Header().Get("X-RateLimit-Limit"))
	}
	if w1.Header().Get("X-RateLimit-Remaining") != "59" {
		t.Fatalf("remaining header %q", w1.Header().Get("X-RateLimit-Remaining"))
	}
	if w1.Header().Get("X-RateLimit-Reset") == "" {
		t.Fatal("missing reset header")
	}

	var last *httptest.ResponseRecorder
	for i := 0; i < 60; i++ {
		last = click("203.0.113.10")
	}
	if last.Code != 429 {
		t.Fatalf("expected 429 on 61st, got %d %s", last.Code, last.Body.String())
	}
	if last.Header().Get("X-RateLimit-Remaining") != "0" {
		t.Fatalf("429 remaining %q", last.Header().Get("X-RateLimit-Remaining"))
	}
	if last.Header().Get("Retry-After") == "" {
		t.Fatal("missing Retry-After")
	}
	var body map[string]any
	if err := json.Unmarshal(last.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["error"] != "rate_limited" {
		t.Fatalf("429 body %#v", body)
	}

	pubw := httptest.NewRecorder()
	h.ServeHTTP(pubw, httptest.NewRequest(http.MethodGet, "/api/v1/public/pages/rl", nil))
	var pub page.PublicPage
	if err := json.Unmarshal(pubw.Body.Bytes(), &pub); err != nil || len(pub.Links) != 1 {
		t.Fatalf("public %#v err=%v body=%s", pub, err, pubw.Body.String())
	}
	if pub.Links[0].Clicks != 60 {
		t.Fatalf("clicks should stay at 60 after 429, got %d", pub.Links[0].Clicks)
	}

	other := click("203.0.113.11")
	if other.Code != 200 {
		t.Fatalf("other ip should be allowed, %d %s", other.Code, other.Body.String())
	}
}
