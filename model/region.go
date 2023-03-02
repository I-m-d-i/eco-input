package model

import (
	"context"
	"database/sql"
	"eco-pasport-input/db"
	"log"
)

type GeneralInformation struct {
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	IsTown bool    `json:"isTown"`
	Oktmo  string  `json:"oktmo"`
}

type EnvironmentalAssessment struct {
	GrossEmissions  string `json:"grossEmissions"`
	WithdrawnWater  string `json:"withdrawnWater"`
	DischargeVolume string `json:"dischargeVolume"`
	FormedWaste     string `json:"formedWaste"`
	Area            string `json:"area"`
	Population      string `json:"population"`
	Center          string `json:"center"`
	CreateDate      int    `json:"createDate"`
}

type Region struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Oktmo string `json:"oktmo"`
}
type RegionInfo struct {
	RegionId                int                     `json:"regionId"`
	Year                    int                     `json:"year"`
	GeneralInformation      GeneralInformation      `json:"generalInformation"`
	EnvironmentalAssessment EnvironmentalAssessment `json:"environmentalAssessment"`
}

func GetRegions() (regions []Region, err error) {
	rows, err := db.Pool.Query(context.Background(), `select r.id, r."name", r.oktmo from eco_pasport."Region" r order by r.sort_region`)
	if err != nil {
		log.Println(err)
		return regions, err
	}
	for rows.Next() {
		var region Region
		if err = rows.Scan(&region.ID, &region.Name, &region.Oktmo); err != nil {
			log.Println(err)
			return regions, err
		}
		regions = append(regions, region)
	}
	return regions, err
}

func SaveRegionInfo(info RegionInfo) error {
	err := saveMunicipalityGeneralInformation(info.GeneralInformation, info.RegionId)
	if err != nil {
		return err
	}
	err = saveMunicipalityEnvironmentalAssessment(info.EnvironmentalAssessment, info.RegionId, info.Year)
	if err != nil {
		return err
	}
	return nil
}

func saveMunicipalityGeneralInformation(generalInformation GeneralInformation, regionId int) error {
	querySaveGeneralInformation := `update "eco-pasport".eco_pasport."Region" set 
		lat=$1,
		lng=$2,
		istown=$3,
		oktmo=$5
	WHERE id = $4`
	_, err := db.Pool.Exec(context.Background(), querySaveGeneralInformation, generalInformation.Lat, generalInformation.Lng, generalInformation.IsTown, regionId, generalInformation.Oktmo)
	if err != nil {
		return err
	}
	return nil
}

func saveMunicipalityEnvironmentalAssessment(environmentalAssessment EnvironmentalAssessment, regionId int, year int) error {
	querySaveEnvironmentalAssessment := `update "eco-pasport".eco_pasport."Region_inform" set discharge_volume=$1,
gross_emissions=$2,
withdrawn_water=$3,
formed_waste=$4,
area=$5,
center=$6,
population_size=$7

							WHERE region_id = $8 AND year = $9`
	_, err := db.Pool.Exec(context.Background(), querySaveEnvironmentalAssessment, environmentalAssessment.DischargeVolume,
		environmentalAssessment.GrossEmissions, environmentalAssessment.WithdrawnWater, environmentalAssessment.FormedWaste,
		environmentalAssessment.Area, environmentalAssessment.Center, environmentalAssessment.Population, regionId, year)
	if err != nil {
		return err
	}
	return nil
}

func GetMunicipalityGeneralInformation(regionId int) (general GeneralInformation, error error) {
	dbRegionInformation := `select r.istown,
	r.lat ,
	r.lng,
	r.oktmo
from eco_pasport."Region_inform" ri
inner join "eco-pasport".eco_pasport."Region" r on r.id = ri.region_id
WHERE ri.region_id = $1`

	information, err := db.Pool.Query(context.Background(), dbRegionInformation, regionId)
	if err != nil {
		log.Println(err)
		return general, err
	}
	for information.Next() {
		var oktmo sql.NullString
		err = information.Scan(&general.IsTown, &general.Lat, &general.Lng, &oktmo)
		general.Oktmo = oktmo.String
		if err != nil {
			log.Println(err)
			return general, err
		}
	}
	return
}

func GetMunicipalityEnvironmentalAssessment(regionId int, year int) (environmentalAssessment EnvironmentalAssessment, error error) {
	dbRegionInformation := `select 
ri.discharge_volume,
ri.gross_emissions,
ri.withdrawn_water,
ri.formed_waste,
ri.area,
ri.center,
ri.population_size,
ri.data_create
							from "eco-pasport".eco_pasport."Region_inform" ri 
							WHERE ri.region_id = $1 AND ri.year = $2`

	information, err := db.Pool.Query(context.Background(), dbRegionInformation, regionId, year)
	if err != nil {
		log.Println(err)
		return environmentalAssessment, err
	}
	for information.Next() {
		if err = information.Scan(&environmentalAssessment.DischargeVolume, &environmentalAssessment.GrossEmissions, &environmentalAssessment.WithdrawnWater,
			&environmentalAssessment.FormedWaste, &environmentalAssessment.Area, &environmentalAssessment.Center, &environmentalAssessment.Population, &environmentalAssessment.CreateDate); err != nil {
			log.Println(err)
			return environmentalAssessment, err
		}
	}
	return environmentalAssessment, nil
}
