package middleware

import (
	"context"
	"forum/internal/middleware/utils"
	"log"
	"net/http"
)

func Handle(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cook, err := r.Cookie("token")
		if err != nil {
			ctx := context.WithValue(r.Context(), "rights", "noauth")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		user, err := utils.ParseToken(r.Context(), cook.Value)
		if err != nil {
			ctx := context.WithValue(r.Context(), "rights", "noauth")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		log.Print(user)

		ctx := context.WithValue(r.Context(), "rights", "auth")

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
