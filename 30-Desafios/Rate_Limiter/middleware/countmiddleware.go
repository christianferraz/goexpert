package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/christianferraz/goexpert/30-Desafios/Rate_Limiter/limiter"
)

func CountMiddleware(next http.HandlerFunc, rateLimiter *limiter.RateLimiter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API_KEY")
		var key string
		if token != "" {
			key = token
		} else {
			key = getIPAddress(r)
		}

		if i, _ := rateLimiter.IsLimited(r.Context(), key); i {

			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}
		if token != "" {
			fmt.Fprintf(w, "Token recebido: %s\n", token)
		} else {
			fmt.Fprintf(w, "IP: %s\n", key)
		}

		next.ServeHTTP(w, r)
	})
}

func getIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-REAL-IP")
	if ip != "" {
		return ip
	}

	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		splitIPs := strings.Split(forwarded, ",")
		return strings.TrimSpace(splitIPs[0])
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}
