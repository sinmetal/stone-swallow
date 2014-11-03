package stone_swallow

import (
	"appengine"
	"net/http"
)

func listBigQuery(w http.ResponseWriter, r *http.Request, c appengine.Context) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
