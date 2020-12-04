package controller

import (
	"net/http"
)

// Static maps static files
func Static(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}
