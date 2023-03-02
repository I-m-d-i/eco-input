package model

import (
	"context"
	"eco-pasport-input/db"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/gbh007/gojlog"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	xlsx "github.com/tealeg/xlsx/v3"
)

// Forms Структура форм приходящих с фронта
type Forms struct {
	Id         int `json:"id"`
	HeaderId   int `json:"headerId"`
	FinalValue struct {
		Text  string `json:"text"`
		Value string `json:"value"`
	} `json:"finalValue"`
	IgnoreAddition bool    `json:"ignoreAddition"`
	OrgNameCol     string  `json:"orgNameCol"`
	IgnoreValue    string  `json:"ignoreValue"`
	LinkedForms    []Forms `json:"linkedForms"`
}

// DataFromJson структура, данных приходящих с фронта
type DataFromJson struct {
	Forms  []Forms `json:"forms"`
	Since  int     `json:"since"`
	Before int     `json:"before"`
	Year   int     `json:"year"`
}

// excelParsOrg структура для таблицы Organization
type excelParsOrg struct {
	year     int
	orgName  string
	orgValue []orgValueDB
	lat      float64
	lng      float64
	regionId int
	address  string
	inn      int
	stroke   int
}

type errorStackStruct struct {
	latLngErrorStack      string
	regionErrorStack      string
	orgNameErrorStack     string
	innErrorStack         string
	bustingFormErrorStack string
	emptyCell             string
}

type regions struct {
	id   int
	name string
}

type orgValueDB struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// ParserExcelOrg Читаем excel с данными организаций для загрузки информации на карту эко-паспорта
func ParserExcelOrg(fileName string, forms []Forms, year, since, before int) error {
	defer duration(track("time"))

	wb, err := xlsx.OpenFile(fileName) // Считываем файл
	if err != nil {
		gojlog.Error(err)
		err = errors.New("Ошибка считывания файла ")
		return err
	}
	sheetName := "Лист1" // Здесь указывается имя листа
	sh, ok := wb.Sheet[sheetName]
	if !ok {
		gojlog.Error(err)
		return errors.New("Такой лист не найден ")
	}

	var region regions
	var regionSlice []regions
	var dataFromExcel []excelParsOrg
	var errs error
	var errorsStruct errorStackStruct
	var wg sync.WaitGroup
	ch := make(chan struct{})
	errorStack := ""

	rows, err := db.Pool.Query(context.Background(), `select 
															r.id ,
															r.name
														  from eco_pasport."Region" r`)
	if err != nil {
		gojlog.Error(err)
	}
	for rows.Next() {
		if err = rows.Scan(&region.id, &region.name); err != nil {
			gojlog.Error(err)
		}
		regionSlice = append(regionSlice, region)
	}

	for i := since - 1; i <= before-1; i++ {
		lineExcel := getDataFromExcel(sh, i, year, &errorsStruct, regionSlice)
		if lineExcel.orgName != "" {
			dataFromExcel = append(dataFromExcel, lineExcel)
		}
	}

	if errorsStruct.latLngErrorStack != "" || errorsStruct.regionErrorStack != "" ||
		errorsStruct.innErrorStack != "" || errorsStruct.emptyCell != "" {
		errorStack +=
			"\n---\nОшибки преобразования долготы и широты\n" + errorsStruct.latLngErrorStack +
				"\n---\nОшибки при поиске идентификатора региона\n" + errorsStruct.regionErrorStack +
				"\n---\nОшибки преобразования инн\n" + errorsStruct.innErrorStack +
				"\n---\nОшибки c пустыми ячейками \n" + errorsStruct.emptyCell
		deleteFile(fileName, err)
		return errors.New(errorStack)
	}

	for i := 0; i <= len(dataFromExcel)-1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createOrgValue(forms, &errorsStruct, &dataFromExcel[i])
		}()
		time.Sleep(5 * time.Millisecond)
	}
	//даём определенное время потокам на выполнение, если они не завершат работу, то вернёт ошибку
	waitTimeout(&wg, ch, &errorStack, before)
	if errorsStruct.bustingFormErrorStack != "" || errorsStruct.orgNameErrorStack != "" || errorStack != "" {
		errorStack += "\n---\nОшибки при проходе формы\n" + errorsStruct.bustingFormErrorStack +
			"\n---\nОшибки при поиске значений у организации\n" + errorsStruct.orgNameErrorStack
		deleteFile(fileName, err)
		return errors.New(errorStack)
	}

	gojlog.Info(len(dataFromExcel), " Кл-во записей в массиве")
	if len(dataFromExcel) == 0 {
		deleteFile(fileName, err)
		return errors.New("вы загружаете пустой файл")
	}
	for i := 0; i < len(dataFromExcel); i++ {
		errs = insertOrganization(dataFromExcel[i], year)
		if errs != nil {
			errorStack += errs.Error() + "\n"
		}
	}
	if errorStack != "" {
		deleteFile(fileName, err)
		return errors.New(errorStack)
	}

	deleteFile(fileName, err)
	gojlog.Info("Загрузка организаций завершена" + errorStack)
	return nil
}

func deleteFile(fileName string, err error) {
	err = os.Remove(fileName)
	if err != nil {
		gojlog.Error(err)
	}
}

// Берем id региона по названию
func getRegionId(name string, regionSlice []regions) (int, error) {
	name = trimValueExcel(name)
	id := 0
	for _, regionId := range regionSlice {
		if regionId.name == name {
			id = regionId.id
		}
	}
	if id == 0 {
		return 0, errors.New("Такого региона не существует '" + name + "'")
	}
	return id, nil
}

func waitTimeout(wg *sync.WaitGroup, ch chan struct{}, errorStack *string, before int) string {
	ch = make(chan struct{})
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	select {
	case <-ch:
	case <-time.After(10 + time.Duration(before)*time.Second):
		gojlog.Error("Горутины зависли")
		*errorStack += "\nВышло время ожидания (Обратитесь к программисту)\n"
	}
	return *errorStack
}

// trimValueExcel
func trimValueExcel(cell string) string {
	r := regexp.MustCompile(`\s+`)
	cell = strings.TrimSpace(cell)
	cell = r.ReplaceAllString(cell, " ")
	return cell
}

func getDataFromExcel(sh *xlsx.Sheet, i, year int, errorStack *errorStackStruct, regionSlice []regions) excelParsOrg {
	var lineExcel excelParsOrg
	var err error
	var cell *xlsx.Cell
	lineExcel.stroke = i + 1

	lineExcel.year = year // Год
	for j := 0; j <= 8; j++ {
		cell, err = sh.Cell(i, j)
		if err != nil {
			gojlog.Error(err)
		}
		if cell.GetStyle().Fill.FgColor == "" || cell.GetStyle().Fill.FgColor == "FF00B0F0" ||
			cell.GetStyle().Fill.FgColor == "FFFFFFFF" {
			return excelParsOrg{}
		}
		if cell.String() == "" {
			errorStack.emptyCell += "\nНашлось пустое значение в строчке " + strconv.Itoa(lineExcel.stroke)
		}
		switch j {
		case 3:
			cell, _ = sh.Cell(i, 3) // Регион
			lineExcel.regionId, err = getRegionId(cell.String(), regionSlice)
			if err != nil {
				errorStack.regionErrorStack += "\nОшибка при поиске идентификатора региона, проверьте правильность его написания, строка " +
					strconv.Itoa(lineExcel.stroke) + " " + err.Error()
			}
		case 4:
			cell, _ = sh.Cell(i, 4) // Название организации
			lineExcel.orgName = trimValueExcel(cell.String())
		case 5:
			cell, _ = sh.Cell(i, 5) // Адрес организации
			lineExcel.address = trimValueExcel(cell.String())
		case 6:
			cell, _ = sh.Cell(i, 6) // ИНН
			lineExcel.inn, err = strconv.Atoi(trimValueExcel(cell.String()))
			if err != nil {
				errorStack.innErrorStack += "\nОшибка преобразования ИНН, строка " + strconv.Itoa(lineExcel.stroke) + " " + err.Error()
			}
		case 7:
			cell, _ = sh.Cell(i, 7) // Широта
			lineExcel.lat, err = cell.Float()
			if err != nil {
				errorStack.latLngErrorStack += "\nОшибка при преобразовании широты, строка " + strconv.Itoa(lineExcel.stroke) + " " + err.Error()
			}
		case 8:
			cell, _ = sh.Cell(i, 8) // Долгота
			lineExcel.lng, err = cell.Float()
			if err != nil {
				errorStack.latLngErrorStack += "\nОшибка при преобразовании долготы, строка " + strconv.Itoa(lineExcel.stroke) + " " + err.Error()
			}
		}
	}
	if errorStack.latLngErrorStack != "" || errorStack.regionErrorStack != "" ||
		errorStack.innErrorStack != "" || errorStack.emptyCell != "" {
		return excelParsOrg{}
	} else {
		return lineExcel
	}
}

// createOrgValue берутся данные из строчки excel раскидывются в нужные поля структуры, собираются и возвращаются ошибки
func createOrgValue(forms []Forms, errorStack *errorStackStruct, lineExcel *excelParsOrg) {
	//Собираем строку orgValue
	if _, returnValue, err := bustingForms(forms, lineExcel.regionId, lineExcel.orgName); err == nil {
		countMinus := 0
		//Если кол-во минусов совпадает с колличеством форм, то выдаётся ошибка
		for _, value := range returnValue {
			if value.Value == "-" {
				countMinus++
			}
		}
		if countMinus == len(forms) {
			errorStack.orgNameErrorStack += "\nНе было найдено значений у организации: " + lineExcel.orgName + " строка " + strconv.Itoa(lineExcel.stroke)
			return
		}
		lineExcel.orgValue = returnValue
	} else {
		errorStack.bustingFormErrorStack += "\nОшибка прохода формы " + err.Error() + " строка " + strconv.Itoa(lineExcel.stroke) + " (ОБРАТИТЕСЬ к Программисту)\n"
		return
	}
}

// bustingForms проход форм
func bustingForms(forms []Forms, regionId int, orgName string) (orgValueDB, []orgValueDB, error) {
	var orgValueNameString = ""
	var orgValueString = ""
	var orgValues []orgValueDB
	var orgValue orgValueDB

	for i, form := range forms {
		orgNameCol, err := strconv.Atoi(form.OrgNameCol)
		if err != nil {
			gojlog.Error(err, " Ошибка при изменении типа данных")
			return orgValue, nil, err
		}
		ignoreValue, err := strconv.Atoi(form.IgnoreValue)
		if err != nil {
			gojlog.Error(err, " Ошибка при изменении типа данных")
			return orgValue, nil, err
		}
		finalValue, err := strconv.Atoi(form.FinalValue.Value)
		if err != nil {
			gojlog.Error(err, " Ошибка при изменении типа данных")
			return orgValue, nil, err
		}
		//Заходим в рекурсию если есть связные формы
		if len(form.LinkedForms) > 0 {
			if tempOrgValue, _, err := bustingForms(form.LinkedForms, regionId, orgName); err == nil {
				orgValue = tempOrgValue
			} else {
				gojlog.Error(err, " Ошибка при обходе рекурсии")
				return orgValue, nil, err
			}
		}
		//Находим все значения у организаций по региону и headerId
		orgValueNameString, orgValueString, err = getValues(finalValue, form.HeaderId, regionId, orgNameCol,
			ignoreValue, orgName, form.FinalValue.Text)
		if err != nil {
			return orgValue, nil, err
		}
		//Если строка не пустая
		if orgValue.Name != "" && orgValue.Value != "" && orgValue.Value != "-" {
			//Value: равен не -
			value, err := strconv.ParseFloat(orgValue.Value, 64)
			if err != nil {
				gojlog.Error(err, " Ошибка при изменении типа данных")
				return orgValue, nil, err
			}
			//Полученное значение не -
			if orgValueString != "-" {
				tempValue, err := strconv.ParseFloat(orgValueString, 64)
				if err != nil {
					gojlog.Error(err, " Ошибка при изменении типа данных")
					return orgValue, nil, err
				}
				//Приплюсовываем значения
				value += tempValue
			}
			//
			orgValue.Name = orgValueNameString
			orgValue.Value = fmt.Sprintf("%.3f", value)
		} else {
			//Собираем элемент по типу "name":"Образование отходов за год, тонн","value":"-"
			orgValue.Name = orgValueNameString
			orgValue.Value = orgValueString
		}
		//Если это форма родитель, то добавляем в новую переменную получившеся значение при завершении обхода
		//и обнуляем orgFinalString
		if form.IgnoreAddition {
			if i < len(forms)-1 {
				orgValues = append(orgValues, orgValue)
				orgValue.Name = ""
				orgValue.Value = ""
			} else {
				//Если это последняя форма, избаляемся от последней запятой
				orgValues = append(orgValues, orgValue)
				orgValue.Name = ""
				orgValue.Value = ""
				return orgValueDB{}, orgValues, nil
			}
		}
	}
	return orgValue, nil, nil
}

// getValues получаем значение из таблицы Row
func getValues(finalValue, headerId, regionId, orgNameCol, ignoreValue int, orgName, finalValueString string) (string, string, error) {
	dbRows := `select r.value[$1] 
			   from eco_pasport."Row" r inner join eco_pasport."Table" t on t.id = r.table_id
               where t.header_id = $2 and t.region_id = $3 and r.value[$4] = $5 
			   and (r.value[$6] != 'ВСЕГО' and r.value[$6] != 'всего' and r.value[$6] != 'газообразные и жидкие, из них:' 
				and r.value[$6] != 'газообразные и жидкие')`
	rows, err := db.Pool.Query(context.Background(), dbRows,
		finalValue, headerId, regionId, orgNameCol, orgName, ignoreValue)
	if err != nil {
		gojlog.Error(err)
		return finalValueString, "", err
	}
	var orgValueElement = ""
	var orgValue = 0.0
	for rows.Next() {
		if err = rows.Scan(&orgValueElement); err != nil {
			gojlog.Error(err)
			return finalValueString, "", err
		} else {
			var tempValue float64
			//форматируем взятое значение из строки и конвертим в флоат
			orgValueElement = strings.Replace(orgValueElement, ",", ".", -1)
			orgValueElement = strings.Replace(orgValueElement, " ", "", -1)
			orgValueElement = strings.Replace(orgValueElement, "-", "", -1)
			if orgValueElement != "" {
				tempValue, err = strconv.ParseFloat(orgValueElement, 64)
				if err != nil {
					gojlog.Error(err, " Ошибка при изменении типа данных ", " headerId ", headerId, " regionId ", regionId, " orgName ", orgName)
					return finalValueString, "", err
				}
			}
			//Приплюсовывем получившееся значение к финальному
			orgValue += tempValue
		}
	}
	//возвращаем получившееся значение
	if orgValueElement == "" {
		return finalValueString, "-", nil
	}
	return finalValueString, fmt.Sprintf("%.3f", orgValue), nil
}

// insertOrganization
func insertOrganization(org excelParsOrg, year int) error {
	value, err := json.Marshal(org.orgValue)
	if err != nil {
		gojlog.Error(err, " regionId ", org.regionId, " orgName ", org.orgName, " строка ", org.stroke)
		return errors.New("Ошибка преобразовании OrgValue в json формат (Обратитесь к программисту)" + strconv.Itoa(org.stroke))
	}
	//по regionId, year и orgName удаляем найденную запись
	dbDeleteOrg := `delete
					from eco_pasport."Organization" ot
					where ot.org_name = $1 and ot.region_id = $2 and ot."year" = $3`
	result, err := db.Pool.Exec(context.Background(), dbDeleteOrg, org.orgName, org.regionId, year)
	if err != nil || result == nil {
		gojlog.Error(err, " строка ", org.stroke-1)
		return errors.New(err.Error() + " Ошибка при удалении организации " + org.orgName + " строка " + strconv.Itoa(org.stroke))
	}
	//insert в базу
	dbInsert := `insert into eco_pasport."Organization"(year, org_name, org_value, lat, lng, region_id, address, inn) values`
	dbInsert = dbInsert + fmt.Sprintf(`(%d,'%s','%s',%g, %g, %d,'%s',%d)`, org.year, org.orgName,
		value, org.lat, org.lng, org.regionId, org.address, org.inn)
	dbInsert = dbInsert + `; `
	result, err = db.Pool.Exec(context.Background(), dbInsert)
	if err != nil || result == nil {
		gojlog.Error(err, " строка ", org.stroke-1)
		return errors.New("Ошибка при загрузке данных в базу, строка " + strconv.Itoa(org.stroke) + " (Обратитесь к программисту)")
	}
	//Смотрим есть ли дубликаты
	dbSelect := `select count(ot.id)
					from eco_pasport."Organization" ot
					where ot.org_name = $1 and ot.region_id = $2 and ot."year" = $3`
	rows, err := db.Pool.Query(context.Background(), dbSelect, org.orgName, org.regionId, year)
	if err != nil {
		gojlog.Error(err)
	}
	var count = 0
	for rows.Next() {
		if err = rows.Scan(&count); err == nil {
			if count > 1 {
				gojlog.Info("В базе существует больше одной записи с такой организацией", " ", org.orgName)
			}
		}
	}
	return nil
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}
