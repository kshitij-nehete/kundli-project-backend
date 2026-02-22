package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/kshitij-nehete/astro-report/internal/response"
	"go.mongodb.org/mongo-driver/mongo"
)

func HealthHandler(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		err := db.Client().Ping(ctx, nil)
		if err != nil {
			response.WriteJSONError(w, http.StatusServiceUnavailable, "database unreachable")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}
}
