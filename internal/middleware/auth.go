package middleware

import (
	"crypto/subtle"
	"database/sql"
	"net/http"
)

func Auth(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				// fallback for file downloads only
				token = r.URL.Query().Get("token")
			}
			if token == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			// strip "Bearer " prefix
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}
			var stored string
			err := db.QueryRow("SELECT value FROM settings WHERE key='api_token'").Scan(&stored)
			if err != nil || stored == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			// constant-time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(token), []byte(stored)) != 1 {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
