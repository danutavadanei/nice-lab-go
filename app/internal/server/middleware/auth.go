package middleware

import (
	"context"
	"github.com/danutavadanei/nice-lab-go/internal/adapters/mysql"
	"net/http"
)

type AuthenticationMiddleware struct {
	tokenUsers map[string]mysql.User
	authRep    *mysql.AuthTokenRepository
}

func NewAuthenticationMiddleware(
	tokenUsers map[string]mysql.User,
	authRep *mysql.AuthTokenRepository,
) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		tokenUsers: tokenUsers,
		authRep:    authRep,
	}
}

func (amw *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")

		if user, found := amw.tokenUsers[token]; found {
			r = r.WithContext(context.WithValue(r.Context(), "user", user))
			next.ServeHTTP(w, r)
			return
		}

		if user, err := amw.authRep.GetUserByAuthToken(r.Context(), token); err == nil {
			amw.tokenUsers[token] = user
			r = r.WithContext(context.WithValue(r.Context(), "user", user))
			next.ServeHTTP(w, r)

			return
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}
