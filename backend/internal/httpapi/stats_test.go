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

func TestMeStatsAndAdminClicks(t *testing.T) {
	authSvc := auth.NewService()
	pages := page.NewService(page.NewMemoryStore())
	h := New(authSvc, pages)

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

	emptyReq := httptest.NewRequest(http.MethodGet, "/api/v1/me/stats", nil)
	emptyReq.Header.Set("Authorization", "Bearer "+login.Token)
	emptyW := httptest.NewRecorder()
	h.ServeHTTP(emptyW, emptyReq)
	if emptyW.Code != 200 {
		t.Fatalf("empty stats %d %s", emptyW.Code, emptyW.Body.String())
	}
	var empty page.MyStats
	if err := json.Unmarshal(emptyW.Body.Bytes(), &empty); err != nil {
		t.Fatal(err)
	}
	if empty.TotalClicks != 0 || empty.Links == nil || len(empty.Links) != 0 {
		t.Fatalf("empty json %#v body=%s", empty, emptyW.Body.String())
	}

	noAuth := httptest.NewRecorder()
	h.ServeHTTP(noAuth, httptest.NewRequest(http.MethodGet, "/api/v1/me/stats", nil))
	if noAuth.Code != 401 {
		t.Fatalf("jwt required %d %s", noAuth.Code, noAuth.Body.String())
	}

	put := httptest.NewRequest(http.MethodPut, "/api/v1/me/page", bytes.NewBufferString(`{"slug":"demo","displayName":"Demo"}`))
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

	click := httptest.NewRequest(http.MethodPost, "/api/v1/public/pages/demo/links/"+link.ID.String()+"/click", nil)
	clickW := httptest.NewRecorder()
	h.ServeHTTP(clickW, click)
	if clickW.Code != 200 {
		t.Fatalf("click %d %s", clickW.Code, clickW.Body.String())
	}

	meReq := httptest.NewRequest(http.MethodGet, "/api/v1/me/stats", nil)
	meReq.Header.Set("Authorization", "Bearer "+login.Token)
	meW := httptest.NewRecorder()
	h.ServeHTTP(meW, meReq)
	if meW.Code != 200 {
		t.Fatalf("me stats %d %s", meW.Code, meW.Body.String())
	}
	var me page.MyStats
	if err := json.Unmarshal(meW.Body.Bytes(), &me); err != nil {
		t.Fatal(err)
	}
	if me.TotalClicks != 1 || len(me.Links) != 1 || me.Links[0].Clicks != 1 || me.Links[0].Title != "Site" || me.Links[0].URL != "https://example.com" {
		t.Fatalf("me stats %#v body=%s", me, meW.Body.String())
	}

	adminReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/stats", nil)
	adminReq.Header.Set("Authorization", "Bearer "+login.Token)
	adminW := httptest.NewRecorder()
	h.ServeHTTP(adminW, adminReq)
	if adminW.Code != 200 {
		t.Fatalf("admin stats %d %s", adminW.Code, adminW.Body.String())
	}
	var admin map[string]any
	if err := json.Unmarshal(adminW.Body.Bytes(), &admin); err != nil {
		t.Fatal(err)
	}
	if admin["users"] != float64(1) || admin["pages"] != float64(1) || admin["totalClicks"] != float64(1) {
		t.Fatalf("admin json %#v body=%s", admin, adminW.Body.String())
	}
	raw, ok := admin["pageClicks"].([]any)
	if !ok || len(raw) != 1 {
		t.Fatalf("pageClicks %#v body=%s", admin["pageClicks"], adminW.Body.String())
	}
	row, ok := raw[0].(map[string]any)
	if !ok {
		t.Fatalf("pageClicks[0] %#v", raw[0])
	}
	if row["clicks"] != float64(1) || row["slug"] == "" || row["id"] == "" {
		t.Fatalf("pageClicks row %#v body=%s", row, adminW.Body.String())
	}
}
