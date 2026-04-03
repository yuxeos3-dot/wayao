package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type loginReq struct {
	Password string `json:"password"`
}

type loginResp struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type changePwReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (app *App) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid request")
		return
	}
	stored := getSetting(app.DB, "api_token")
	hash := sha256Hash(req.Password)

	if stored == "" {
		// first login: default password is "admin"
		if req.Password != "admin" {
			jsonError(w, 401, "invalid password")
			return
		}
		// generate and store a new token
		token := generateToken()
		app.DB.Exec("UPDATE settings SET value=? WHERE key='api_token'", token)
		// also store the password hash for future logins
		app.DB.Exec("INSERT OR REPLACE INTO settings(key,value) VALUES('password_hash',?)", hash)
		jsonOK(w, loginResp{
			Token:     token,
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		})
		return
	}

	// normal login: verify password hash
	storedHash := getSetting(app.DB, "password_hash")
	if storedHash == "" {
		// legacy: no password_hash, token is the password
		if hash != stored && req.Password != stored {
			jsonError(w, 401, "invalid password")
			return
		}
	} else {
		if hash != storedHash {
			jsonError(w, 401, "invalid password")
			return
		}
	}

	jsonOK(w, loginResp{
		Token:     stored,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
	})
}

func (app *App) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	var req changePwReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid request")
		return
	}
	if req.NewPassword == "" || len(req.NewPassword) < 6 {
		jsonError(w, 400, "password must be at least 6 characters")
		return
	}

	// verify old password
	storedHash := getSetting(app.DB, "password_hash")
	oldHash := sha256Hash(req.OldPassword)
	if storedHash != "" && oldHash != storedHash {
		jsonError(w, 401, "old password incorrect")
		return
	}

	newHash := sha256Hash(req.NewPassword)
	newToken := generateToken()
	app.DB.Exec("UPDATE settings SET value=? WHERE key='password_hash'", newHash)
	app.DB.Exec("UPDATE settings SET value=? WHERE key='api_token'", newToken)

	jsonOK(w, loginResp{
		Token:     newToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
	})
}

func sha256Hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func generateToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("crypto/rand failed: %v", err))
	}
	return hex.EncodeToString(b)
}
