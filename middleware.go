package main

import (
	"log"
	"net/http"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
)

// Auth checks whether request token is valid
func (a *App) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// leave user header empty for no user authenticated
		r.Header.Add("User", "")

		// retrieve and parse jwt
		jwt, err := jws.ParseJWTFromRequest(r)
		// no token found in request
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// do a lookup of token in blacklist
		serializedToken := string(r.Header.Get("Authorization")[7:])
		_, err = a.Redis.Get(serializedToken).Result()

		// token is found in blacklist so do not continue
		if err == nil {
			next.ServeHTTP(w, r)
			return
		}

		// check if token is valid
		if err := jwt.Validate(a.RSAKey.Public(), crypto.SigningMethodRS256); err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// check if email is found in claim
		if jwt.Claims().Get("id") != nil {
			// insert email of user from token to request header
			r.Header.Set("User", jwt.Claims().Get("id").(string))
		}

		// serve route endpoint with user in header
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs user request for endpoint
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s > %s %s", r.RemoteAddr, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
