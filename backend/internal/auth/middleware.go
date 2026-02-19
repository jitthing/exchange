package auth

import (
	"net/http"
	"strings"
)

const defaultDemoUser = "demo-user"

// extractBearerToken extracts a Bearer token from the Authorization header.
func extractBearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return ""
}

// RequireAuth rejects with 401 if no valid token is present.
// When AUTH_DISABLED=true, injects "demo-user" and continues.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsAuthDisabled() {
			ctx := contextWithUserID(r.Context(), defaultDemoUser)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		token := extractBearerToken(r)
		if token == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		userID, err := ValidateToken(token)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		ctx := contextWithUserID(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth sets user if token present, continues either way.
func OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsAuthDisabled() {
			ctx := contextWithUserID(r.Context(), defaultDemoUser)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		token := extractBearerToken(r)
		if token != "" {
			if userID, err := ValidateToken(token); err == nil {
				ctx := contextWithUserID(r.Context(), userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
