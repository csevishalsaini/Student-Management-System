package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu        sync.Mutex
	visitor   map[string]int
	limit     int
	resetTime time.Duration
}

func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter {
	r1 := &rateLimiter{
		visitor:   make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}
	go r1.resetVisitorCount()
	return r1
}

func (r1 *rateLimiter) resetVisitorCount() {
	for {
		time.Sleep(r1.resetTime)
		r1.mu.Lock()
		r1.visitor = make(map[string]int)
		r1.mu.Unlock()
	}
}

func (r1 *rateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ✅ Extract only IP (not port)
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		r1.mu.Lock()
		r1.visitor[ip]++
		count := r1.visitor[ip]
		r1.mu.Unlock()

		fmt.Printf("Visitor count from %v is %v\n", ip, count)

		// ✅ If limit exceeded → STOP execution
		if count > r1.limit {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"Too many requests"}`))
			return // 🔥 VERY IMPORTANT
		}

		next.ServeHTTP(w, r)
	})
}