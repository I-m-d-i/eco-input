package rows

import (
	"eco-pasport-input/model"
	"encoding/json"
	"gitlab.com/gbh007/gojlog"
	"log"
	"net/http"
)

// Сохраняем новые данные с формы в БД
func saveValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	type DataFromFront struct {
		Region int      `json:"id_region"`
		Header int      `json:"id_header"`
		Value  []string `json:"value"`
	}
	var jsonData DataFromFront
	var receivedRow model.Row
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	receivedRow.Value = jsonData.Value
	var tableId = model.GetTableID(jsonData.Region, jsonData.Header)
	err = model.AddValue(receivedRow, tableId)
	if err != nil {
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	log.Println("Запись добавлена\nТаблица: ", tableId, "Значения: ", receivedRow.Value)
}

// Редактирование данных в базе
func editRow(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	type DataFromFront struct {
		Id    int      `json:"id"`
		Value []string `json:"value"`
	}
	var jsonData DataFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var editableRow = model.Row{Id: jsonData.Id, Value: jsonData.Value}
	errOfModel := model.EditRow(editableRow)
	if errOfModel != nil {
		log.Println(errOfModel)
		http.Error(w, "", http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(errOfModel.Error()); err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}
	log.Println("Запись отредактирована\nСтрока ", editableRow.Id, ", значения: ", editableRow.Value)
}

// Дублирование информации в базе
func copyRow(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	var copiedRow model.Row
	err := json.NewDecoder(r.Body).Decode(&copiedRow)
	if err != nil {
		log.Println(err)
	}
	var tableId = copiedRow.TableId()
	errOfModel := model.AddValue(copiedRow, tableId)
	if errOfModel != nil {
		log.Println(errOfModel, "\nТаблица: ", tableId, "Значения: ", copiedRow.Value)
		http.Error(w, "", http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(errOfModel.Error()); err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}
	log.Println("Запись скопированная \nТаблица: ", tableId, "Значения: ", copiedRow.Value)
}

// Меняем местами строки в таблице
func changeSort(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	var elementSort []model.Row
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	err := json.NewDecoder(r.Body).Decode(&elementSort)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	model.ChangeSortRows(elementSort)
}

// Удаление строки
func deleteRow(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	var row model.Row
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&row.Id); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := model.DeleteRow(row); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Удалена строчка: ", row.Id)
}

// Переносим строчку с таблицы в другой регион
func transferRegion(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//TODO можно отказаться от IdHeader и IdRegion так как эту информаццию можно получить по id строки
	type IDFromFront struct {
		IdHeader       int `json:"id_header"`
		IdRegion       int `json:"id_region"`
		IdRegionTarget int `json:"id_region_target"`
		IdRow          int `json:"id_row"`
	}
	var jsonData IDFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var row model.Row
	row.Id = jsonData.IdRow
	errOfModel := model.TransferRowInRegion(row, model.GetTableID(jsonData.IdRegionTarget, jsonData.IdHeader))
	if errOfModel != nil {
		log.Println(errOfModel)
		http.Error(w, "", http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(errOfModel.Error()); err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}
	log.Println("Строка: ", row.Id, " из таблицы ", jsonData.IdHeader, "\n Перенесена из ", jsonData.IdRegion, " в ", jsonData.IdRegionTarget)
}

// Получаем строчки с БД
func getValues(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	type IDFromFront struct {
		IdHeader int `json:"id_header"`
		IdRegion int `json:"id_region"`
	}
	var jsonData IDFromFront
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tableID := model.GetTableID(jsonData.IdRegion, jsonData.IdHeader)
	var response = model.GetRows(tableID)
	// Отправляем данные на фронт
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func calculateFinalValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST, OPTIONS")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var fromFront model.CalculateForm
	err := json.NewDecoder(r.Body).Decode(&fromFront)
	if err != nil {
		gojlog.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = model.CalculateTotalValue(fromFront)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
