package tables

import (
	"eco-pasport-input/model"
	"encoding/json"
	"log"
	"net/http"
)

// Переименование таблицы
func renameTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	type DataFromFront struct {
		TableId int    `json:"id"`
		NewName string `json:"newName"`
	}
	var jsonData DataFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = model.RenameTable(jsonData.TableId, jsonData.NewName)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	type DataFromFront struct {
		Year   int                `json:"year"`
		Name   string             `json:"name"`
		Header []model.JsonHeader `json:"header"`
		Groups []string           `json:"groups"`
	}
	var jsonData DataFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	idTable, err := model.AddTable(jsonData.Year, jsonData.Name, jsonData.Groups)
	if err != nil {
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}
	err = model.EditHeader(jsonData.Header, idTable)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func deleteTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	type DataFromFront struct {
		TableId int `json:"id"`
	}
	var jsonData DataFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = model.DeleteTable(jsonData.TableId)
	if err != nil {
		log.Println(err)
		if err := json.NewEncoder(w).Encode("В таблицу внесены данные, удаление невозможно"); err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}
}

func saveHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	type DataFromFront struct {
		TableId int                `json:"id"`
		Header  []model.JsonHeader `json:"header"`
	}
	var jsonData DataFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = model.EditHeader(jsonData.Header, jsonData.TableId)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func saveSortTables(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	type DataFromFront struct {
		Tables []model.Table `json:"tables"`
	}
	var jsonData DataFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = model.EditSort(jsonData.Tables)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func saveGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	type DataFromFront struct {
		Groups   []string `json:"groups"`
		IdHeader int      `json:"idHeader"`
	}
	var jsonData DataFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = model.SaveTableGroups(jsonData.Groups, jsonData.IdHeader)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Запрос на получение списка всех таблиц
func getTables(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}

	// Получение данных с фронта
	w.Header().Set("Content-Type", "application/json")

	type IDFromFront struct {
		Year int `json:"year"`
	}
	var jsonYear IDFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonYear)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := model.GetTables(jsonYear.Year)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
