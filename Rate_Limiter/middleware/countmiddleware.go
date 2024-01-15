package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func CountMiddleware(next http.HandlerFunc, ctx *context.Context, rdb *redis.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := fmt.Sprintf("%s-%s", getIPAdress(r), r.URL.Path)
		val, err := rdb.Get(*ctx, key).Result()
		if err != nil {
			if err != redis.Nil {
				panic(err)
			}
			val = "0"
		}
		count, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		if count > 10 {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		count++
		w.Write([]byte(fmt.Sprintf("Você é o visitante número %v\n", getIPAdress(r))))
		fmt.Println("eu sou o key", getIPAdress(r))
		if err := rdb.Set(*ctx, key, count, 10*time.Second).Err(); err != nil {
			panic(err)
		}
		next.ServeHTTP(w, r)
	})
}

func getIPAdress(r *http.Request) string {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}
	return ""
}
