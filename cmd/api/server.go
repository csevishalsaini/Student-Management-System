package main

import (
	"crypto/tls"
	"os"

	"fmt"
	"log"
	"net/http"
	mw "restapi/internal/api/middleware"
	"restapi/internal/api/router"
	"restapi/internal/repository/sqlconnect"

	"github.com/joho/godotenv"
)

func main() {

		err := godotenv.Load()
		if(err != nil){
			fmt.Println(err)
			return
		}
		_,err =sqlconnect.ConnectDb()
		if(err !=nil){
			fmt.Println("Error------ ",err)
			return
		}

	key := "key.pem"
	cert := "cert.pem"

	port := os.Getenv("API_PORT")

	router := router.Router()

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// r1 := mw.NewRateLimiter(5, time.Minute)

	// hppOtions := mw.HPPOptions{
	// 	CheckQuery:                  true,
	// 	CheckBody:                   true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	Whitelist:                   []string{"name", "sortby", "sortorder"},
	// }

	// securemux := utils.ApplyMiddlewares(mux, mw.Hpp((hppOtions)), mw.Compression, mw.SecurityHeaders, mw.ResponseTimeMiddleware, r1.Middleware)
	// securemux = mw.SecurityHeaders(mux)
	securemux := mw.SecurityHeaders(router)

	server := &http.Server{
		Addr: port,
		// Handler: middlewares.SecurityHeaders(mux),
		// Handler: middlewares.Cors(mux),
		// Handler: mw.Cors(r1.Middleware(mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Compression(mw.Hpp(hppOtions)(mux)))))),
		Handler:   securemux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port", port)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatal("Error starting the server", err)
	}

}
