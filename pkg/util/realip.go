package util

import (
	"net"
	"net/http"
	"strings"
)

func RealIP(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	if remoteAddr == "" {
		return ""
	}
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}

	ip := net.ParseIP(host)
	if !ip.IsLoopback() {
		return host
	}

	xff := r.Header.Get("X-Forwarded-For")
	if xff == "" {
		return host
	}

	for _, part := range strings.Split(xff, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if hp, _, err := net.SplitHostPort(part); err == nil {
			part = hp
		}
		return part
	}

	return host
}
