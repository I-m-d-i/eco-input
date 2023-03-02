package model

import (
	"context"
	"database/sql"
	"eco-pasport-input/db"
	"errors"
	"fmt"
	"log"
	"strconv"
)

type Table struct {
	Id     int      `json:"id"`
	Name   string   `json:"name"`
	Sort   int      `json:"sort"`
	Groups []string `json:"groups"`
}

// GetTableID Берем id таблиц по региону и хэдеру TODO: Возвращать и обрабатывать ошибку
func GetTableID(idRegion int, idHeader int) int {
	rows, err := db.Pool.Query(context.Background(), `select t.id from eco_pasport."Table" t where t.header_id = $1 and t.region_id = $2`,
		idHeader, idRegion)
	if err != nil {
		log.Println(err)
		return 0
	}
	var idTable int
	for rows.Next() {
		if err := rows.Scan(&idTable); err != nil {
			log.Println(err)
			return 0
		}
	}
	return idTable
}

// TODO: Возвращать и обрабатывать ошибку
func GetTables(year int) []Table {
	rows, err := db.Pool.Query(context.Background(), `select h.id,  h."name", h.sort, h."permissionCode" from eco_pasport."Header" h where h.year = $1 order by h.sort`, year)
	if err != nil {
		log.Println(err)
	}
	var tables []Table
	for rows.Next() {
		var table Table
		if err = rows.Scan(&table.Id, &table.Name, &table.Sort, &table.Groups); err != nil {
			log.Println(err)
		}
		tables = append(tables, table)
	}
	return tables
}

func RenameTable(tableId int, name string) (err error) {
	name = trimValue(name)
	isNotDub, err := checkNameTable(tableId, name, 0)
	if err != nil {
		return
	}
	if !isNotDub {
		err = errors.New("таблица с таким именем уже есть")
		return
	}
	request := `update eco_pasport."Header" set name = $1 where id = $2`
	_, err = db.Pool.Exec(context.Background(), request, name, tableId)
	if err != nil {
		return
	}
	return
}

func AddTable(year int, name string, groups []string) (int, error) {
	var sort sql.NullInt32
	var idTable int
	name = trimValue(name)
	isNotDub, err := checkNameTable(0, name, year)
	if err != nil {
		return 0, err
	}
	if !isNotDub {
		err := errors.New("таблица с таким именем уже есть")
		return 0, err
	}
	requestMaxSort := `SELECT max(h.sort) FROM eco_pasport."Header" h  WHERE "year"= $1`
	if err = db.Pool.QueryRow(context.Background(), requestMaxSort, year).Scan(&sort); err != nil {
		log.Println(err)
		return 0, err
	}
	request := `insert into eco_pasport."Header" (name,"year" ,sort,"permissionCode") values($1,$2 ,$3,$4) RETURNING id`
	if err = db.Pool.QueryRow(context.Background(), request, name, year, sort.Int32+1, groups).Scan(&idTable); err != nil {
		log.Println(err)
		return 0, err
	}
	requestCross := `	
	Insert into eco_pasport."Table" (header_id,region_id)
	SELECT 
	h.id AS header_id, 
	r.id AS region_id
	FROM "eco-pasport".eco_pasport."Header" h CROSS JOIN eco_pasport."Region" r
	WHERE h.id = $1`
	_, err = db.Pool.Exec(context.Background(), requestCross, idTable)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return idTable, nil
}

func DeleteTable(tableId int) (err error) {
	request := `DELETE FROM eco_pasport."Header" WHERE id = $1`
	_, err = db.Pool.Exec(context.Background(), request, tableId)
	if err != nil {
		return fmt.Errorf("В таблицу " + strconv.Itoa(tableId) + " внесены данные, удаление невозможно")
	}
	return
}

func EditSort(tables []Table) (err error) {
	for i := 0; i < len(tables); i++ {
		request := `update eco_pasport."Header" set sort = $1 where id = $2`
		_, err = db.Pool.Exec(context.Background(), request, tables[i].Sort, tables[i].Id)
		if err != nil {
			return
		}
	}
	return
}

func CloningTables(year int, cloningYear int) (err error) {
	request := `INSERT INTO eco_pasport."Header" (name,"header","year",sort,"permissionCode")
	SELECT
		name,
		"header" ,
		$1,
		sort,
		"permissionCode"
	FROM eco_pasport."Header" 
	WHERE year = $2`
	_, err = db.Pool.Exec(context.Background(), request, year, cloningYear)
	if err != nil {
		return
	}
	requestCross := `	
	Insert into eco_pasport."Table" (header_id,region_id)
	SELECT 
	h.id AS header_id, 
	r.id AS region_id
	FROM "eco-pasport".eco_pasport."Header" h CROSS JOIN eco_pasport."Region" r
	WHERE h."year"  = $1`
	_, err = db.Pool.Exec(context.Background(), requestCross, year)
	if err != nil {
		return
	}
	return
}

func SaveTableGroups(groups []string, idHeader int) error {
	request := `update eco_pasport."Header" set "permissionCode" = $1 where id = $2`
	_, err := db.Pool.Exec(context.Background(), request, groups, idHeader)
	if err != nil {
		return err
	}
	return nil
}

func checkNameTable(tableId int, name string, year int) (bool, error) {
	var countElement int // Счетчик одинаковых элементов в таблице
	if year == 0 {
		requestYear := `
		SELECT h.year
		FROM "eco-pasport".eco_pasport."Header" h
		WHERE h.id = $1`
		resultYear, err := db.Pool.Query(context.Background(), requestYear, tableId)
		if err != nil {
			return true, err
		}
		for resultYear.Next() {
			if err = resultYear.Scan(&year); err != nil {
				return true, err
			}
		}
	}
	request := `
	SELECT count(h.id)
	FROM eco_pasport."Header" h
	WHERE h.name = $2 AND h.year=$1`
	resultRequest, err := db.Pool.Query(context.Background(), request, year, name)
	if err != nil {
		return true, err
	}
	for resultRequest.Next() {
		if err = resultRequest.Scan(&countElement); err != nil {
			return true, err
		}
	}
	return countElement <= 0, nil
}
