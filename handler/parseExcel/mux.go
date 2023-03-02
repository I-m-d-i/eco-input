package parseExcel

import (
	"eco-pasport-input/model"
	"encoding/json"
	"gitlab.com/gbh007/gojlog"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func ApplyMux(mux *http.ServeMux) {
	mux.HandleFunc("/api/uploadFile", getExcel)
	mux.HandleFunc("/api/createFile", createTemplate)
	mux.HandleFunc("/api/getTemplate", getTemplate)
}

var handler *multipart.FileHeader
var file multipart.File

// Получаем excel с Front
func getExcel(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	w.Header().Set("Content-Type", "multipart/form-data")
	w.Header().Set("Content-Type", "application/json")

	contentType := strings.Split(r.Header.Get("Content-Type"), ";")
	var err error

	if contentType[0] == "multipart/form-data" {
		if err := r.ParseMultipartForm(1024); err != nil {
			gojlog.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// get 'file'
		file, handler, err = r.FormFile("file")
		if err != nil {
			gojlog.Error(err)
			err := file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// read file bytes
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			gojlog.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// write bytes to a localfile
		err = os.WriteFile(handler.Filename, fileBytes, 0644)
		if err != nil {
			gojlog.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if contentType[0] == "application/json" {
		var fromFront model.DataFromJson
		err := json.NewDecoder(r.Body).Decode(&fromFront)
		if err != nil {
			gojlog.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = model.ParserExcelOrg(handler.Filename, fromFront.Forms, fromFront.Year, fromFront.Since, fromFront.Before)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Получаем данные с БД
func createTemplate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	type DataFromFront struct {
		Data string `json:"data"`
		Year int    `json:"year"`
	}
	var fromFront DataFromFront
	err := json.NewDecoder(r.Body).Decode(&fromFront)
	if err != nil {
		gojlog.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = model.CreateTemplate(fromFront.Data, fromFront.Year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTemplate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "GET")

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		return
	}
	w.Header().Set("Content-Type", "application/json")

	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		gojlog.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	returnString, err := model.GetTemplate(year)
	if err != nil {
		if err.Error() == "open configs/"+strconv.Itoa(year)+".json: The system cannot find the file specified." {
			http.Error(w, "4452", http.StatusInternalServerError)
			return
		}
		if err.Error() == "open configs/"+strconv.Itoa(year)+".json: no such file or directory" {
			http.Error(w, "4452", http.StatusInternalServerError)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(returnString); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
