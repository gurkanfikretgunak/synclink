package httpapi

import (
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"strconv"
	"strings"
	"sync"
	"time"
)

type clickLimiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	hits   map[string][]time.Time
}

func newClickLimiter(limit int, window time.Duration) *clickLimiter {
	if limit < 1 {
		limit = 60
	}
	if window <= 0 {
		window = time.Minute
	}
	return &clickLimiter{limit: limit, window: window, hits: map[string][]time.Time{}}
}

func (l *clickLimiter) allow(key string, now time.Time) (ok bool, remaining int, reset time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	cut := now.Add(-l.window)
	prev := l.hits[key]
	kept := prev[:0]
	for _, t := range prev {
		if t.After(cut) {
			kept = append(kept, t)
		}
	}
	reset = now.Add(l.window)
	if len(kept) > 0 {
		reset = kept[0].Add(l.window)
	}
	if len(kept) >= l.limit {
		l.hits[key] = kept
		return false, 0, reset
	}
	kept = append(kept, now)
	l.hits[key] = kept
	return true, l.limit - len(kept), reset
}

func clientIP(r *http.Request) string {
	if x := r.Header.Get("X-Forwarded-For"); x != "" {
		return strings.TrimSpace(strings.Split(x, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func setRateHeaders(w http.ResponseWriter, limit, remaining int, reset time.Time) {
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
	w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(reset.UTC().Unix(), 10))
}

func (s *Server) limitClicks(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := clientIP(r) + "|" + chi.URLParam(r, "id")
		ok, remaining, reset := s.clicks.allow(key, time.Now())
		setRateHeaders(w, s.clicks.limit, remaining, reset)
		if !ok {
			w.Header().Set("Retry-After", strconv.Itoa(60))
			writeJSON(w, 429, map[string]any{"error": "rate_limited", "message": "too many clicks", "code": 429})
			return
		}
		next.ServeHTTP(w, r)
	})
}
