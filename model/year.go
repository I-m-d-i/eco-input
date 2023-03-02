package model

import (
	"context"
	"eco-pasport-input/db"
	"log"
	"sort"
)

func GetYears() ([]int, error) {
	var years []int
	rows, err := db.Pool.Query(context.Background(), `select h.year from eco_pasport."Header" h group by year`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var year int
		if err = rows.Scan(&year); err != nil {
			log.Println(err)
			return nil, err
		}
		years = append(years, year)
	}
	sort.SliceStable(years, func(i, j int) bool { return years[i] > years[j] })
	return years, err
}

func GetYearsTree() ([]int, error) {
	var years []int
	rows, err := db.Pool.Query(context.Background(), `select T."Year" from eco_pasport."Tree" T group by "Year"`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var year int
		if err = rows.Scan(&year); err != nil {
			log.Println(err)
			return nil, err
		}
		years = append(years, year)
	}
	sort.SliceStable(years, func(i, j int) bool { return years[i] > years[j] })
	return years, err
}

func GetYearsRegions(regionId int) ([]int, error) {
	var years []int
	rows, err := db.Pool.Query(context.Background(), `select  ri.year from eco_pasport."Region_inform" ri  where ri.region_id=$1 order by ri."year" desc`, regionId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var year int
		if err = rows.Scan(&year); err != nil {
			log.Println(err)
			return nil, err
		}
		years = append(years, year)
	}
	//sort.SliceStable(years, func(i, j int) bool { return years[i] > years[j] })
	return years, err
}
