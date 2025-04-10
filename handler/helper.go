package handler

import (
	"context"
	"net/http"
	"time"
)

func getContextFromRequest(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), 5*time.Minute)
}
