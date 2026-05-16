package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		token := parts[1]
		user, err := s.validateToken(token)
		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", fmt.Sprintf("%d", user.ID))
		r.Header.Set("X-Username", user.Username)

		next(w, r)
	}
}

func (s *Server) validateToken(token string) (*AuthUser, error) {
	// Format: base64(username:password)
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid token format")
	}

	username := parts[0]
	return &AuthUser{
		ID:       1,
		Username: username,
		Role:     "admin",
	}, nil
}

func generateToken(username string) string {
	token := fmt.Sprintf("%s:authenticated", username)
	return base64.StdEncoding.EncodeToString([]byte(token))
}

// validateCredentials checks if username/password are valid
func (s *Server) validateCredentials(username, password string) bool {
	// currently hardcoded.
	return username == "admin" && password == "password"
}
