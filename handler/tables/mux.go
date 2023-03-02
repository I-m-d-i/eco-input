package tables

import "net/http"

func ApplyMux(mux *http.ServeMux) {
	//TODO Преобразовать в вид api/table/*
	mux.HandleFunc("/api/renameTable", renameTable)
	mux.HandleFunc("/api/addTable", addTable)
	mux.HandleFunc("/api/deleteTable", deleteTable)
	mux.HandleFunc("/api/saveHeader", saveHeader)
	mux.HandleFunc("/api/saveSortTables", saveSortTables)
	mux.HandleFunc("/api/table", getTables)
	mux.HandleFunc("/api/table/saveGroups", saveGroups)
}
