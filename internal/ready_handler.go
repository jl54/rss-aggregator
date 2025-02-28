package internal

import "net/http"

func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, 200, struct{}{})
}
