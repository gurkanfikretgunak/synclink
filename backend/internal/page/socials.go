package page

import (
	"net/mail"
	"net/url"
	"strings"
)

const maxSocials = 12

type Social struct {
	Network string `json:"network"`
	URL     string `json:"url"`
}

var allowedNetworks = map[string]struct{}{
	"instagram": {},
	"x":         {},
	"youtube":   {},
	"tiktok":    {},
	"github":    {},
	"linkedin":  {},
	"threads":   {},
	"spotify":   {},
	"website":   {},
	"email":     {},
}

func emptySocials() []Social {
	return []Social{}
}

func copySocials(in []Social) []Social {
	if len(in) == 0 {
		return emptySocials()
	}
	out := make([]Social, len(in))
	copy(out, in)
	return out
}

func NormalizeSocials(in []Social) []Social {
	if len(in) == 0 {
		return emptySocials()
	}
	out := make([]Social, 0, len(in))
	for _, s := range in {
		if len(out) >= maxSocials {
			break
		}
		n := strings.ToLower(strings.TrimSpace(s.Network))
		if n == "twitter" {
			n = "x"
		}
		if _, ok := allowedNetworks[n]; !ok {
			continue
		}
		u, ok := normalizeSocialURL(n, s.URL)
		if !ok {
			continue
		}
		out = append(out, Social{Network: n, URL: u})
	}
	return out
}

func normalizeSocialURL(network, raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", false
	}
	if network == "email" {
		return normalizeEmailURL(raw)
	}
	return normalizeHTTPURL(raw)
}

func normalizeHTTPURL(raw string) (string, bool) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", false
	}
	if u.Host == "" {
		return "", false
	}
	return raw, true
}

func normalizeEmailURL(raw string) (string, bool) {
	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, "mailto:") {
		addr := strings.TrimSpace(raw[len("mailto:"):])
		if validEmailAddr(addr) {
			return "mailto:" + addr, true
		}
		return "", false
	}
	if validEmailAddr(raw) {
		return raw, true
	}
	return normalizeHTTPURL(raw)
}

func validEmailAddr(addr string) bool {
	if addr == "" || strings.ContainsAny(addr, " \t\n\r") {
		return false
	}
	_, err := mail.ParseAddress(addr)
	return err == nil
}
