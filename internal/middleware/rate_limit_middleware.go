package middleware

import (
	"net"
	"net/http"
	"sync"

	"github.com/kshitij-nehete/astro-report/internal/handler"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter *rate.Limiter
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(1, 5) // 1 req/sec, burst 5
		visitors[ip] = &visitor{limiter}
		return limiter
	}

	return v.limiter
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		limiter := getVisitor(ip)

		if !limiter.Allow() {
			handler.WriteJSONError(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}

		next.ServeHTTP(w, r)
	})
}
