package middlewares

import (
	"net/http"
	"strings"
)

func MiddlewareExcludePaths(middlewares func(http.Handler) http.Handler, excluded ...string) func(http.Handler) http.Handler{
 return func(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _,path:= range excluded{
				if strings.HasPrefix(r.URL.Path,path){
					next.ServeHTTP(w,r)
					return
				}
		}
		middlewares(next).ServeHTTP(w,r)
	})
 }

}
