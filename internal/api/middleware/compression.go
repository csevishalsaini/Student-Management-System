package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

func Compression(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if(!strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")){
			next.ServeHTTP(w,r)
		}
		
		w.Header().Set("Content-Encoding","gzip")
		gzip := gzip.NewWriter(w)
		defer gzip.Close()

		w = &gzipResponseWriter{ResponseWriter: w, Writer: gzip}
		next.ServeHTTP(w,r)
		fmt.Println("Sent response from Compression Middleware")
		

	})
}

type gzipResponseWriter struct{
	http.ResponseWriter
	Writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(b [] byte)(int,error){
	return g.Writer.Write(b)
}