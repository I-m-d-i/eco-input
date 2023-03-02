package model

import (
	"context"
	"errors"
	"log"
	"regexp"
	"strings"

	"eco-pasport-input/db"
	//"github.com/lib/pq"
)

// Row Структура данных для редактирования/копирования/перемещения
type Row struct {
	Sort  int      `json:"sort"`
	Id    int      `json:"id"`
	Value []string `json:"value"`
}

// TableId получаем id таблицы в которой находится строка
func (row *Row) TableId() (TableId int) {
	rows, err := db.Pool.Query(context.Background(), "select r.table_id from eco_pasport.\"Row\" r where r.id = $1", row.Id)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		if err := rows.Scan(&TableId); err != nil {
			log.Println(err)
			return
		}
	}
	return TableId
}
func DeleteRow(row Row) (err error) {
	_, err = db.Pool.Exec(context.Background(), "delete from eco_pasport.\"Row\" r where r.id = $1", row.Id)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

// GetRows Получение строк для фронта из выбранной таблицы
func GetRows(IdTable int) (Rows []Row) {
	rowsDb, err := db.Pool.Query(context.Background(), `select r.id, sort, r.value from eco_pasport."Row" r
                           where r.table_id = $1 order by sort desc`, IdTable)
	if err != nil {
		log.Println(err)
		return
	}
	for rowsDb.Next() {
		var rowTemp Row
		if err := rowsDb.Scan(&rowTemp.Id, &rowTemp.Sort, &rowTemp.Value); err != nil {
			log.Println(err)
			return
		}
		Rows = append(Rows, rowTemp)
	}
	return Rows
}

// AddValue Добавление строки в указанную таблицу
func AddValue(addedRow Row, tableId int) error {
	for i, value := range addedRow.Value {
		addedRow.Value[i] = trimValue(value)
	}
	if checkValue(addedRow.Value, tableId) {
		request := `insert into eco_pasport."Row"(value , table_id ) values ($1,$2)`
		result, err := db.Pool.Exec(context.Background(), request, addedRow.Value, tableId)
		if err != nil || result == nil {
			return err
		}
		return err
	} else {
		err := errors.New("запись уже существует")
		return err
	}
}

// EditRow Изменение указанной строки
func EditRow(editableRow Row) error {
	for i, value := range editableRow.Value {
		editableRow.Value[i] = trimValue(value)
	}
	if checkValue(editableRow.Value, editableRow.TableId()) {
		request := `update eco_pasport."Row" set value = $1 where id = $2`
		result, err := db.Pool.Exec(context.Background(), request, editableRow.Value, editableRow.Id)
		if err != nil || result == nil {
			return err
		}
		return nil
	} else {
		err := errors.New("запись уже существует")
		return err
	}
}

// ChangeSortRows Запрос на замену строк местами
func ChangeSortRows(Rows []Row) {
	requestChangeSort(Rows[0].Sort, Rows[1].Id)
	requestChangeSort(Rows[1].Sort, Rows[0].Id)
	log.Println("Поменяли местами строки: ", "Id:", Rows[0].Id, " Sort:", Rows[0].Sort, " и ", "Id:", Rows[1].Id, " Sort:", Rows[1].Sort)
}

// Изменение сортировки
func requestChangeSort(whereTransf int, whatTransf int) {
	_, err := db.Pool.Exec(context.Background(), `update eco_pasport."Row" set sort =$1 where id = $2`, whereTransf, whatTransf)
	if err != nil {
		log.Println(err)
		return
	}
}

// Проверка вводимой информации на дубли
func checkValue(Value []string, TableID int) bool {
	var countElement int // Счетчик одинаковых элементов в таблице
	request := `select count(r.id) from eco_pasport."Row" r where r.table_id = $1 and r.value = $2`
	resultRequest, err := db.Pool.Query(context.Background(), request, TableID, Value)
	if err != nil {
		log.Println(err)
	}
	for resultRequest.Next() {
		if err = resultRequest.Scan(&countElement); err != nil {
			log.Println(err)
		}
	}
	return countElement <= 0
}

// Избавляемся от перевода строки и лишних пробелов
func trimValue(Value string) string {
	re := regexp.MustCompile(`\n`)
	r := regexp.MustCompile(`\s+`)
	Value = strings.TrimSpace(Value)
	Value = re.ReplaceAllString(Value, " ")
	Value = r.ReplaceAllString(Value, " ")
	return Value
}

// TransferRowInRegion  перенос строки в другой регион
// TODO: необходимо передавать id целевого региона, а не таблицу
func TransferRowInRegion(portableRow Row, idTableTarget int) (err error) {
	resultRequest, err := db.Pool.Query(context.Background(), `select value from eco_pasport."Row" where id = $1`, portableRow.Id)
	if err != nil {
		return
	}
	var portableRowValue []string
	for resultRequest.Next() {
		if err = resultRequest.Scan(&portableRowValue); err != nil {
			return
		}
	}
	if checkValue(portableRowValue, idTableTarget) {
		_, err = db.Pool.Exec(context.Background(), `update eco_pasport."Row" set table_id = $1 where id = $2`,
			idTableTarget, portableRow.Id)
		if err != nil {
			return
		}
	} else {
		err = errors.New("запись там уже существует, перенос отменен")
		return
	}
	return
}
