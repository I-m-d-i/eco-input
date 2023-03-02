package rows

import "net/http"

func ApplyMux(mux *http.ServeMux) {
	mux.HandleFunc("/api/getRows", getValues)
	mux.HandleFunc("/api/saveRow", saveValue)
	mux.HandleFunc("/api/deleteRow", deleteRow)
	mux.HandleFunc("/api/editRow", editRow)
	mux.HandleFunc("/api/copyRow", copyRow)
	mux.HandleFunc("/api/changeSortRows", changeSort)
	mux.HandleFunc("/api/transferRowRegion", transferRegion)
	mux.HandleFunc("/api/calculateFinalValue", calculateFinalValue)
}
