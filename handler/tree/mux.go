package tree

import "net/http"

func ApplyMux(mux *http.ServeMux) {
	mux.HandleFunc("/api/getTree", getTree)
	mux.HandleFunc("/api/yearsTree", getYearsTree)
	mux.HandleFunc("/api/addTree", addTree)
}
