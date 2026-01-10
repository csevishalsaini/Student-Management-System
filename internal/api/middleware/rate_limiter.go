package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct{
		mu sync.Mutex
		visitor map[string]int
		limit int
		resetTime time.Duration
}

func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter{
		r1 := &rateLimiter{visitor: make(map[string]int), limit:  limit, resetTime: resetTime}
		go r1.resetVisitorCount()
		return r1
}

func (r1 *rateLimiter) resetVisitorCount(){
	for{
		time.Sleep(r1.resetTime)
		r1.mu.Lock()
		r1.visitor = make(map[string]int)
		r1.mu.Unlock()

	}
}

func (r1 *rateLimiter)Middleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r1.mu.Lock()
			defer r1.mu.Unlock()
			
			visitorIP := r.RemoteAddr;
			r1.visitor[visitorIP]++
			fmt.Printf("Visitor count from %v is %v ",visitorIP,r1.visitor[visitorIP])
			if(r1.visitor[visitorIP]>r1.limit){
				http.Error(w,"Too many request", http.StatusTooManyRequests)
			}
			
			next.ServeHTTP(w,r)

	})
}