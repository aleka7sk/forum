package post

import "net/http"

func Error404(w http.ResponseWriter) {
	http.Error(w, "404 Source not Found", 404)
}

func Error505(w http.ResponseWriter) {
	http.Error(w, "500 Internal Server Error", 500)
}
