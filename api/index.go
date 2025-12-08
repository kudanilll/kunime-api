package api

import (
	"kunime-api/pkg/app"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	app.Handler(w, r)
}
