package model

import (
	"context"
	"eco-pasport-input/db"
	"fmt"
	"gitlab.com/gbh007/gojlog"
	"strconv"
	"strings"
	"unicode"
)

type CalculateForm struct {
	RegionId         int    `json:"regionId"`
	HeaderId         int    `json:"headerId"`
	DesiredTextValue string `json:"desiredTextValue"`
	InputText        string `json:"inputText"`
	SearchIn         string `json:"searchIn"`
	IgnoreValue      string `json:"ignoreValue"`
}

func CalculateTotalValue(form CalculateForm) error {
	var tableId = GetTableID(form.RegionId, form.HeaderId)

	dbRows := `select r.value
			   from eco_pasport."Row" r
			   where r.table_id = $1 and r.value[$2+1] = $3`
	rows, err := db.Pool.Query(context.Background(), dbRows, tableId, form.SearchIn, form.DesiredTextValue)
	if err != nil {
		gojlog.Error(err)
		return err
	}
	var row Row
	var totalValue Row

	for rows.Next() {
		err = rows.Scan(&row.Value)
		if err != nil {
			gojlog.Error(err)
			return err
		}
		for index, value := range row.Value {
			value = strings.Replace(value, ",", ".", -1)
			value = strings.Replace(value, " ", "", -1)
			if len(totalValue.Value) == 0 {
				totalValue.Value = make([]string, len(row.Value))
			}
			if isLetter(value) || index == 0 {
				totalValue.Value[index] += ""
				continue
			}
			tempTotal, _ := strconv.ParseFloat(totalValue.Value[index], 64)
			tempValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				gojlog.Error(err)
				return err
			}
			totalValue.Value[index] = fmt.Sprintf("%0.3f", tempTotal+tempValue)
		}
	}

	index, _ := strconv.Atoi(form.SearchIn)
	totalValue.Value[index] = form.InputText
	if form.IgnoreValue != "" {
		index, _ = strconv.Atoi(form.IgnoreValue)
		totalValue.Value[index] = ""
	}

	for i, value := range totalValue.Value {
		totalValue.Value[i] = strings.Replace(value, ".", ",", -1)
	}

	err = AddValue(totalValue, tableId)
	if err != nil {
		return err
	}
	return nil
}

func isLetter(value string) bool {
	if len(value) == 0 {
		return true
	}
	if value == "-" || value == "" {
		return true
	}
	for _, letter := range value {
		if unicode.IsLetter(letter) == true {
			return true
		}
	}
	return false
}
