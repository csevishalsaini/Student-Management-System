package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"restapi/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

func JwtMiddleware(next http.Handler) http.Handler {
	fmt.Println("------------------ JWT Middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("++++++++++++++++ Inside JWT Middleware")
		token, err := r.Cookie("Bearer")

		if err != nil {
			http.Error(w, "Authorization Header Missing ", http.StatusUnauthorized)
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")

		parsedToken, err := jwt.Parse(token.Value, func(token *jwt.Token) (any, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v ", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "Token expired ", http.StatusUnauthorized)
				return
			}else if errors.Is(err, jwt.ErrTokenMalformed){
				http.Error(w, "Token Malformed ", http.StatusUnauthorized)
				return
			}
			utils.ErrorHandler(err, "")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return 
		}


		if parsedToken.Valid {
			fmt.Println("Valid jwt")
		}else{
			http.Error(w, "Invalid", http.StatusUnauthorized)
			log.Println("Invalid jwt", token.Value)
			return
		}
	
 		claims, ok := parsedToken.Claims.(jwt.MapClaims);
		if !ok {
			http.Error(w, "Invalid login TOKEN", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "role", claims["role"])
		ctx = context.WithValue(r.Context(), "expiresAt", claims["exp"])
		ctx = context.WithValue(r.Context(), "username", claims["user"])
		ctx = context.WithValue(r.Context(), "userId", claims["uid"])

		fmt.Println(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))

		fmt.Println("++++++++++++++++ Sent Response from JWT Middleware")
	})
}



/*
























*/