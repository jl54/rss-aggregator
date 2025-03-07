package internal

import "net/http"

func HandleError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Something went wrong")
}
