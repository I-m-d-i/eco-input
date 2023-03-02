package regions

import (
	"eco-pasport-input/model"
	"encoding/json"
	"log"
	"net/http"
)

func ApplyMux(mux *http.ServeMux) {
	mux.HandleFunc("/api/regions", getRegions)
	mux.HandleFunc("/api/yearsRegions", yearsRegion)
	mux.HandleFunc("/api/getGeneralInformation", getGeneralInformation)
	mux.HandleFunc("/api/getEnvironmentalAssessment", getMunicipalityEnvironmentalAssessment)
	mux.HandleFunc("/api/saveRegionInfo", saveRegionInfo)
}

// Запрос на получение списка регионов
func getRegions(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	regions, err := model.GetRegions()
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(regions); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func yearsRegion(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	type IDFromFront struct {
		RegionId int `json:"regionId"`
	}
	var jsonData IDFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	w.Header().Set("Content-Type", "application/json")
	years, err := model.GetYearsRegions(jsonData.RegionId)
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

func saveRegionInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	var regionInfo model.RegionInfo
	err := json.NewDecoder(r.Body).Decode(&regionInfo)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = model.SaveRegionInfo(regionInfo)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		if err = json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	log.Println("Информация для МО ", regionInfo.RegionId, " за ", regionInfo.Year, " год была обновлена")
}

func getGeneralInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	type IDFromFront struct {
		RegionId int `json:"regionId"`
	}
	var jsonData IDFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	years, err := model.GetMunicipalityGeneralInformation(jsonData.RegionId)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(years); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func getMunicipalityEnvironmentalAssessment(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	type IDFromFront struct {
		RegionId int `json:"regionId"`
		Year     int `json:"year"`
	}
	var jsonData IDFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	years, err := model.GetMunicipalityEnvironmentalAssessment(jsonData.RegionId, jsonData.Year)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(years); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
