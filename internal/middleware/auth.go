package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pratham15541/go-crud/internal/models"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				sendAuthError(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Extract token from "Bearer <token>" format
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				sendAuthError(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := tokenParts[1]

			// Parse and validate token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				sendAuthError(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Token is valid, continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// sendAuthError sends an authentication error response
func sendAuthError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResp := models.ErrorResponse{
		Error:   "Authentication Error",
		Message: message,
		Code:    statusCode,
	}

	json.NewEncoder(w).Encode(errorResp)
}