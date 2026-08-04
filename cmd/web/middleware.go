package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v5"
)

func LoadSession(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// redirectToLogin manda al login recordando a dónde iba el usuario, para
// devolverlo ahí tras autenticarse (parámetro "next").
func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	login := "/login?next=" + url.QueryEscape(r.URL.RequestURI())
	http.Redirect(w, r, login, http.StatusSeeOther)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := session.GetString(r.Context(), "token")

		if tokenString == "" {
			log.Println("Not authenticated")
			redirectToLogin(w, r)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Secreto compartido con el servicio auth, cargado de .env en run().
			return jwtSecret, nil
		})
		if err != nil {
			log.Println("Error parsing token: ", err)
			redirectToLogin(w, r)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			log.Println("Error reading claims", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// create new context from `r` request context, and assign key `"user"`
		// to value of `"123"`
		ctx := context.WithValue(r.Context(), "user", claims["user"])

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new `"user"` key
		// will be accessible from this point forward.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
