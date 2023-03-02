package model

import (
	"context"
	"eco-pasport-input/db"
	"encoding/json"
	"log"
)

type JsonHeader struct {
	Text string       `json:"text"`
	Sub  []JsonHeader `json:"sub"`
}

// GetHeader Получение шапки по idHeader
func GetHeader(idHeader int) (headerGo []JsonHeader, err error) {
	var headerBase []byte
	err = db.Pool.QueryRow(context.Background(), `select  h."header" from eco_pasport."Header" h where h.id = $1`, idHeader).Scan(&headerBase)
	if err != nil {
		log.Println(err)
		return headerGo, err
	}

	if err = json.Unmarshal(headerBase, &headerGo); err != nil {
		log.Println(err)
		return headerGo, nil
	}
	return headerGo, nil
}

// EditHeader Измяет шапку паблицы по idHeader
func EditHeader(header []JsonHeader, idHeader int) error {
	headerByte, err := json.Marshal(header)
	if err != nil {
		return err
	}
	_, err = db.Pool.Exec(context.Background(), `UPDATE eco_pasport."Header" SET header = $1 where id = $2`, headerByte, idHeader)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
