package model

import (
	"os"
	"strconv"
)

func CreateTemplate(data string, year int) error {
	var jsonBlob = []byte(data)
	err := os.WriteFile("configs/"+strconv.Itoa(year)+".json", jsonBlob, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetTemplate(year int) (string, error) {
	b, err := os.ReadFile("configs/" + strconv.Itoa(year) + ".json")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
