package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func ResponseTimeMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Recived Request in Response Time")
		start := time.Now()
		
		wrappedWriter := &responseWriter{ResponseWriter: w, code:http.StatusOK}
		duration := time.Since(start)
		w.Header().Set("X-Response-Time",duration.String())
		next.ServeHTTP(wrappedWriter,r)

		duration = time.Since(start)
		fmt.Printf("Method %s, URL: %s, status: %d, Duration: %v \n",r.Method,r.URL,wrappedWriter.code, duration)
		fmt.Println("Sent request from Response Time middleware ")


	})
}

type responseWriter struct{
	http.ResponseWriter
	code int
}

func (rw *responseWriter ) WriteHeader(code int){
	rw.code = code
	rw.ResponseWriter.WriteHeader(code)
}