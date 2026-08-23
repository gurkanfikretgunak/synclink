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

func registerAndPage(t *testing.T, h http.Handler) string {
	t.Helper()
	reg := httptest.NewRecorder()
	h.ServeHTTP(reg, httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(`{"email":"a@b.com","password":"password1"}`)))
	if reg.Code != 201 {
		t.Fatalf("register %d %s", reg.Code, reg.Body.String())
	}
	var login struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(reg.Body.Bytes(), &login); err != nil || login.Token == "" {
		t.Fatalf("token %#v", login)
	}
	put := httptest.NewRequest(http.MethodPut, "/api/v1/me/page", bytes.NewBufferString(`{"slug":"demo","displayName":"Demo"}`))
	put.Header.Set("Authorization", "Bearer "+login.Token)
	put.Header.Set("Content-Type", "application/json")
	pw := httptest.NewRecorder()
	h.ServeHTTP(pw, put)
	if pw.Code != 200 {
		t.Fatalf("page %d %s", pw.Code, pw.Body.String())
	}
	return login.Token
}

func TestPublicSubscribeAndMeSubscribers(t *testing.T) {
	h := New(auth.NewService(), page.NewService(page.NewMemoryStore()))
	token := registerAndPage(t, h)

	bad := httptest.NewRecorder()
	h.ServeHTTP(bad, httptest.NewRequest(http.MethodPost, "/api/v1/public/pages/demo/subscribe", bytes.NewBufferString(`{"email":"nope"}`)))
	if bad.Code != 400 {
		t.Fatalf("invalid %d %s", bad.Code, bad.Body.String())
	}
	missing := httptest.NewRecorder()
	h.ServeHTTP(missing, httptest.NewRequest(http.MethodPost, "/api/v1/public/pages/nope/subscribe", bytes.NewBufferString(`{"email":"ok@ex.com"}`)))
	if missing.Code != 404 {
		t.Fatalf("missing %d %s", missing.Code, missing.Body.String())
	}

	ok := httptest.NewRecorder()
	h.ServeHTTP(ok, httptest.NewRequest(http.MethodPost, "/api/v1/public/pages/demo/subscribe", bytes.NewBufferString(`{"email":"Fan@Ex.com"}`)))
	if ok.Code != 201 {
		t.Fatalf("sub %d %s", ok.Code, ok.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(ok.Body.Bytes(), &body); err != nil || body["ok"] != true {
		t.Fatalf("json %#v err=%v", body, err)
	}
	dup := httptest.NewRecorder()
	h.ServeHTTP(dup, httptest.NewRequest(http.MethodPost, "/api/v1/public/pages/demo/subscribe", bytes.NewBufferString(`{"email":"fan@ex.com"}`)))
	if dup.Code != 409 {
		t.Fatalf("dup %d %s", dup.Code, dup.Body.String())
	}

	list := httptest.NewRequest(http.MethodGet, "/api/v1/me/subscribers", nil)
	list.Header.Set("Authorization", "Bearer "+token)
	lw := httptest.NewRecorder()
	h.ServeHTTP(lw, list)
	if lw.Code != 200 {
		t.Fatalf("list %d %s", lw.Code, lw.Body.String())
	}
	var subs []page.SubscriberDTO
	if err := json.Unmarshal(lw.Body.Bytes(), &subs); err != nil || len(subs) != 1 || subs[0].Email != "fan@ex.com" {
		t.Fatalf("subs %#v err=%v body=%s", subs, err, lw.Body.String())
	}

	del := httptest.NewRequest(http.MethodDelete, "/api/v1/me/subscribers/"+subs[0].ID.String(), nil)
	del.Header.Set("Authorization", "Bearer "+token)
	dw := httptest.NewRecorder()
	h.ServeHTTP(dw, del)
	if dw.Code != 204 {
		t.Fatalf("delete %d %s", dw.Code, dw.Body.String())
	}
}

func TestPublicPagePasswordLocked(t *testing.T) {
	h := New(auth.NewService(), page.NewService(page.NewMemoryStore()))
	token := registerAndPage(t, h)
	put := httptest.NewRequest(http.MethodPut, "/api/v1/me/page", bytes.NewBufferString(`{"slug":"demo","displayName":"Demo","pagePassword":"secret"}`))
	put.Header.Set("Authorization", "Bearer "+token)
	put.Header.Set("Content-Type", "application/json")
	pw := httptest.NewRecorder()
	h.ServeHTTP(pw, put)
	if pw.Code != 200 {
		t.Fatalf("lock page %d %s", pw.Code, pw.Body.String())
	}

	locked := httptest.NewRecorder()
	h.ServeHTTP(locked, httptest.NewRequest(http.MethodGet, "/api/v1/public/pages/demo", nil))
	if locked.Code != 401 {
		t.Fatalf("locked %d %s", locked.Code, locked.Body.String())
	}
	var body map[string]any
	_ = json.Unmarshal(locked.Body.Bytes(), &body)
	if body["error"] != "locked" {
		t.Fatalf("locked json %#v", body)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/pages/demo", nil)
	req.Header.Set("X-Page-Password", "secret")
	ok := httptest.NewRecorder()
	h.ServeHTTP(ok, req)
	if ok.Code != 200 {
		t.Fatalf("unlock %d %s", ok.Code, ok.Body.String())
	}
	var pub page.PublicPage
	if err := json.Unmarshal(ok.Body.Bytes(), &pub); err != nil || pub.Slug != "demo" {
		t.Fatalf("pub %#v err=%v", pub, err)
	}
}

func TestLinkExtrasJSON(t *testing.T) {
	h := New(auth.NewService(), page.NewService(page.NewMemoryStore()))
	token := registerAndPage(t, h)
	cr := httptest.NewRequest(http.MethodPost, "/api/v1/me/page/links", bytes.NewBufferString(`{"title":"Site","url":"https://example.com","featured":true,"sensitive":true,"thumbnailUrl":"https://img.example/a.png"}`))
	cr.Header.Set("Authorization", "Bearer "+token)
	cr.Header.Set("Content-Type", "application/json")
	cw := httptest.NewRecorder()
	h.ServeHTTP(cw, cr)
	if cw.Code != 201 {
		t.Fatalf("link %d %s", cw.Code, cw.Body.String())
	}
	var raw map[string]any
	if err := json.Unmarshal(cw.Body.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	if raw["featured"] != true || raw["sensitive"] != true || raw["thumbnailUrl"] != "https://img.example/a.png" {
		t.Fatalf("link json %#v", raw)
	}
	pubw := httptest.NewRecorder()
	h.ServeHTTP(pubw, httptest.NewRequest(http.MethodGet, "/api/v1/public/pages/demo", nil))
	var pub page.PublicPage
	if err := json.Unmarshal(pubw.Body.Bytes(), &pub); err != nil || len(pub.Links) != 1 || !pub.Links[0].Featured || !pub.Links[0].Sensitive {
		t.Fatalf("public extras %#v err=%v body=%s", pub, err, pubw.Body.String())
	}
}
