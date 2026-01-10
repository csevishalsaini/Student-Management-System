package middlewares

import "net/http"

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Disable DNS prefetching
		w.Header().Set("X-DNS-Prefetch-Control", "off")
		
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// XSS protection (legacy browsers)
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Prevent MIME-type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Enforce HTTPS
		w.Header().Set(
			"Strict-Transport-Security",
			"max-age=63072000; includeSubDomains; preload",
		)

		// Content Security Policy
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'",
		)

		// Referrer policy
		w.Header().Set("Referrer-Policy", "no-referrer")

		// Hide server info
		w.Header().Set("X-Powered-By", "")

		w.Header().Set("Server", "")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")

		// Cross-Origin isolation headers
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

		// CORS headers
		
		// Browser permissions
		w.Header().Set("Permissions-Policy", "geolocation=(self), microphone=()")

		next.ServeHTTP(w, r)
	})
}

