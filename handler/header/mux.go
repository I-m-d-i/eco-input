package header

import (
	"eco-pasport-input/model"
	"encoding/json"
	"log"
	"net/http"
)

func ApplyMux(mux *http.ServeMux) {
	mux.HandleFunc("/api/getHeader", getHeader)
}

// Запрос на получение заголовков таблиц
func getHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	type DataFromFront struct {
		Header int `json:"id_header"`
	}
	var jsonDate DataFromFront
	// Получение данных с фронта
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&jsonDate)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Берем json с базы
	headerGo, err := model.GetHeader(jsonDate.Header)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Отправляем данные на фронт
	if err = json.NewEncoder(w).Encode(headerGo); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
