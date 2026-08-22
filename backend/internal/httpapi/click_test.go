package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gurkanfikretgunak/synclink/backend/internal/auth"
	"github.com/gurkanfikretgunak/synclink/backend/internal/page"
)

func TestPublicClickIncrementsAnd404(t *testing.T) {
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

	click := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/public/pages/demo/links/"+link.ID.String()+"/click", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		return w
	}
	w1 := click()
	if w1.Code != 200 {
		t.Fatalf("click1 %d %s", w1.Code, w1.Body.String())
	}
	var body1 map[string]any
	if err := json.Unmarshal(w1.Body.Bytes(), &body1); err != nil {
		t.Fatal(err)
	}
	if body1["ok"] != true || body1["clicks"] != float64(1) {
		t.Fatalf("click json %#v", body1)
	}
	w2 := click()
	var body2 map[string]any
	_ = json.Unmarshal(w2.Body.Bytes(), &body2)
	if w2.Code != 200 || body2["clicks"] != float64(2) {
		t.Fatalf("click2 %#v %s", body2, w2.Body.String())
	}

	pubw := httptest.NewRecorder()
	h.ServeHTTP(pubw, httptest.NewRequest(http.MethodGet, "/api/v1/public/pages/demo", nil))
	var pub page.PublicPage
	if err := json.Unmarshal(pubw.Body.Bytes(), &pub); err != nil || len(pub.Links) != 1 || pub.Links[0].Clicks != 2 {
		t.Fatalf("public %#v err=%v body=%s", pub, err, pubw.Body.String())
	}

	bad := httptest.NewRecorder()
	h.ServeHTTP(bad, httptest.NewRequest(http.MethodPost, "/api/v1/public/pages/demo/links/"+uuid.New().String()+"/click", nil))
	if bad.Code != 404 {
		t.Fatalf("missing link %d %s", bad.Code, bad.Body.String())
	}
	missingPage := httptest.NewRecorder()
	h.ServeHTTP(missingPage, httptest.NewRequest(http.MethodPost, "/api/v1/public/pages/nope/links/"+link.ID.String()+"/click", nil))
	if missingPage.Code != 404 {
		t.Fatalf("missing page %d %s", missingPage.Code, missingPage.Body.String())
	}
	inactive := false
	pr := httptest.NewRequest(http.MethodPatch, "/api/v1/me/page/links/"+link.ID.String(), bytes.NewBufferString(`{"active":false}`))
	pr.Header.Set("Authorization", "Bearer "+login.Token)
	pr.Header.Set("Content-Type", "application/json")
	pw2 := httptest.NewRecorder()
	h.ServeHTTP(pw2, pr)
	if pw2.Code != 200 {
		t.Fatalf("hide %d %s", pw2.Code, pw2.Body.String())
	}
	hid := httptest.NewRecorder()
	h.ServeHTTP(hid, httptest.NewRequest(http.MethodPost, "/api/v1/public/pages/demo/links/"+link.ID.String()+"/click", nil))
	if hid.Code != 404 {
		t.Fatalf("inactive %d %s", hid.Code, hid.Body.String())
	}
	_ = inactive
}
