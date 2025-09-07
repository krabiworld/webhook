package middlewares

import (
	"errors"
	"net/http"
	"strings"
	"webhook/internal/db"
	"webhook/internal/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		_, err := db.G[models.Token]().Where("token = ?", token).First(r.Context())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Error().Err(err).Msg("Authentication failed")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
