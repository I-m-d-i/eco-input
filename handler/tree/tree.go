package tree

import (
	"eco-pasport-input/model"
	"encoding/json"
	"log"
	"net/http"
)

// Добавление сруктуры папок
func addTree(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	type DataFromFront struct {
		Year int        `json:"year"`
		Tree model.Tree `json:"tree"`
	}
	var jsonData DataFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = jsonData.Tree.AddInDB(jsonData.Year)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Внесены изменение в дерево экопаспортов за", jsonData.Year, " год")
}

func getTree(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	type DataFromFront struct {
		Year int `json:"year"`
	}
	var year DataFromFront
	if err := json.NewDecoder(r.Body).Decode(&year); err != nil && err.Error() != "EOF" {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := model.GetTree(year.Year)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

}

func getYearsTree(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	years, err := model.GetYearsTree()
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
