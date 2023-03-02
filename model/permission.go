package model

import (
	"context"
	"eco-pasport-input/db"
	"log"
)

type Group struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// DeletePermissionGroup Удаляет группу пользователей
func DeletePermissionGroup(group Group) error {
	_, err := db.Pool.Exec(context.Background(), `delete from eco_pasport."User_groups" where id=$1`, group.Id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// SavePermissionGroup Применяет изменения в данные группы пользователей
func SavePermissionGroup(group Group) error {
	_, err := db.Pool.Exec(context.Background(), `update eco_pasport."User_groups" set code = $1, name = $2 where id=$3`, group.Code, group.Name, group.Id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// AddPermissionGroup добавляет новую группу пользователей
func AddPermissionGroup(group Group) error {
	_, err := db.Pool.Exec(context.Background(), `insert into eco_pasport."User_groups" (name,code) values($1,$2)`, group.Name, group.Code)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetPermissionGroups получает все группы пользователей
// TODO дописать медод сохранения
func GetPermissionGroups() (groups []Group, err error) {
	rows, err := db.Pool.Query(context.Background(), `select id, code, name from eco_pasport."User_groups"`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var group Group
	for rows.Next() {
		if err = rows.Scan(&group.Id, &group.Code, &group.Name); err != nil {
			log.Println(err)
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}
