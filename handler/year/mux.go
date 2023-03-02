package year

import (
	"eco-pasport-input/model"
	"encoding/json"
	"log"
	"net/http"
)

func ApplyMux(mux *http.ServeMux) {
	mux.HandleFunc("/api/addYear", cloningTables)
	mux.HandleFunc("/api/years", getYears)
}

// Получаем список всех годов
func getYears(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	years, err := model.GetYears()
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(years); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// TODO: перенести в таблицы
func cloningTables(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	type DataFromFront struct {
		Year        int `json:"year"`
		CloningYear int `json:"cloningYear"`
	}
	var jsonData DataFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = model.CloningTables(jsonData.Year, jsonData.CloningYear)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
