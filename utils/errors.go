package utils

import (
	"net/http"
)

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	http.Redirect(w, r, "/oops", 500)

	// w.Write([]byte("Internal server error"))

}

func UnauthorizedUserError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Need authorization"))
}
