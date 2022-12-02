package authorization

import (
	"github.com/cristalhq/jwt/v4"
	"log"
	"os"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

var verifier, _ = jwt.NewVerifierHS(jwt.HS256, jwtKey)

//func Authorizer(e *casbin.Enforcer) func(next http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//
//		fn := func(w http.ResponseWriter, r *http.Request) {
//
//			token, err := jwt.Parse([]byte(r.Header.Get("token")), verifier)
//			if err != nil {
//				log.Println(err)
//				http.Error(w, "unauthorized", http.StatusUnauthorized)
//				return
//			}
//
//			claims := GetMapClaims(token.Bytes())
//
//			res, err := e.EnforceSafe(claims["userType"], r.URL.Path, r.Method)
//			if err != nil {
//				log.Println("enforce error")
//				http.Error(w, "unauthorized user", http.StatusUnauthorized)
//				return
//			}
//			log.Println(res)
//
//			if res {
//				log.Println("redirect")
//				next.ServeHTTP(w, r)
//			} else {
//				http.Error(w, "forbidden", http.StatusForbidden)
//				return
//			}
//
//		}
//
//		return http.HandlerFunc(fn)
//	}
//}

func GetToken(tokenString string) *jwt.Token {
	token, err := jwt.Parse([]byte(tokenString), verifier)
	if err != nil {
		log.Println(err)
	}
	return token
}

func GetMapClaims(tokenBytes []byte) map[string]string {
	var claims map[string]string

	err := jwt.ParseClaims(tokenBytes, verifier, &claims)
	if err != nil {
		log.Println(err)
	}

	return claims
}
