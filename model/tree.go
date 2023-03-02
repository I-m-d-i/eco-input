package model

import (
	context "context"
	"database/sql"
	"eco-pasport-input/db"
	"errors"
	"github.com/jackc/pgx/v4"
	"log"
	"sort"
	"strconv"
)

// Tree структура для создания дерева на основе данных из TreeDB
type Tree struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Lvl      int    `json:"lvl"`
	Children []Tree `json:"children"`
	IdHeader []int  `json:"idHeader"`
	Sort     int    `json:"sort"`
}

// TreeDB структура для получения данных с базы
type TreeDB struct {
	Id         int           `json:"id"`
	Name       string        `json:"name"`
	Lvl        int           `json:"lvl"`
	IdChildren sql.NullInt32 `json:"children"`
	IdHeader   sql.NullInt32 `json:"idHeader"`
	Sort       int           `json:"sort"`
}

type Leaf struct {
	HeaderId int    `json:"int"`
	Name     string `json:"name"`
	Sort     int    `json:"sort"`
}

var (
	treeDB []TreeDB
	trees  []Tree
	leafs  []Leaf
)

// GetTree формирует дерево из полученных с БД данных по определённому year
func GetTree(year int) (completeTree []Tree, error error) {
	dbTreeRequest := `select
						t."Id" as "Id",
						t."Name",
						t."Lvl" ,
						t2."Id" as "Child_Id",
						h."id" as "Id_Header",
						t."sort"
						from "eco-pasport".eco_pasport."Tree" t
						left join "eco-pasport".eco_pasport."Tree" t2 on t2."Id_Parent"=t."Id"
						left join "eco-pasport".eco_pasport."Tree_Header" th on t."Id"=th."Id_Tree"
						left join "eco-pasport".eco_pasport."Header" h on h.id = th."Id_Header"
						where t."Year" = $1
						order by "sort" `

	dbLeafsRequest := `select
						t."id",
						t."name",
						t."sort"
						from "eco-pasport".eco_pasport."Header" t where "year" = $1`
	trees = nil
	leafs = nil
	treeDB = nil
	var nodeDB TreeDB
	var leaf Leaf

	leafsRequest, err := db.Pool.Query(context.Background(), dbLeafsRequest, year)
	if err != nil {
		log.Println(err)
		return completeTree, err
	}
	for leafsRequest.Next() {
		if err = leafsRequest.Scan(&leaf.HeaderId, &leaf.Name, &leaf.Sort); err != nil {
			log.Println(err)
			return completeTree, err
		}
		leafs = append(leafs, leaf)
	}
	treeRequest, err := db.Pool.Query(context.Background(), dbTreeRequest, year)
	if err != nil {
		log.Println(err)
		return completeTree, err
	}
	for treeRequest.Next() {
		if err = treeRequest.Scan(&nodeDB.Id, &nodeDB.Name, &nodeDB.Lvl, &nodeDB.IdChildren, &nodeDB.IdHeader, &nodeDB.Sort); err != nil {
			log.Println(err)
			return completeTree, err
		}
		treeDB = append(treeDB, nodeDB)
	}
	initialNodes := make([]int, 0)
	for _, value := range treeDB {
		if value.Lvl == 0 && !Contains(initialNodes, value.Id) {
			initialNodes = append(initialNodes, value.Id)
		}
	}
	for _, id := range initialNodes {
		tree, err := genTree(id)
		if err != nil {
			return trees, err
		}
		trees = append(trees, tree)
	}
	trees = sortTree(trees)
	return trees, nil
}

// Contains ищем уникальные id для 0 lvl после добавляем их в слайл uniq
func Contains(arr []int, key int) bool {
	for _, value := range arr {
		if value == key {
			return true
		}
	}
	return false
}

// genTree строит отдельное дерево для каждого заголовка
func genTree(id int) (Tree, error) {
	var node Tree
	for _, value := range treeDB {
		if value.Id == id {
			node.Id = value.Id
			node.Lvl = value.Lvl
			node.Name = value.Name
			node.Sort = value.Sort
			if value.IdChildren.Valid == true {
				children, err := genTree(int(value.IdChildren.Int32))
				if err != nil {
					return Tree{}, err
				}
				if children.Name != "" {
					node.Children = append(node.Children, children)
				}
			} else if value.IdHeader.Valid == true {
				leaf, err := getLeaf(int(value.IdHeader.Int32), value.Id)
				if err != nil {
					log.Println(err)
					return Tree{}, err
				}
				node.Children = append(node.Children, leaf)
			} else {
				node.Children = make([]Tree, 0)
			}
		}
	}
	return node, nil
}

// getLeaf находит название таблицы по хэдер id и возвращает объект tree в качестве leaf
func getLeaf(IdHeader int, idParent int) (Tree, error) {
	for _, leaf := range leafs {
		if leaf.HeaderId == IdHeader {
			id, _ := strconv.Atoi(strconv.Itoa(IdHeader) + strconv.Itoa(idParent))
			return Tree{
				Id:       id,
				Name:     leaf.Name,
				Sort:     leaf.Sort,
				IdHeader: []int{IdHeader},
			}, nil
		}
	}
	return Tree{}, errors.New("Таблица " + strconv.Itoa(IdHeader) + " была не найдена")
}

// sortTree сортирует существующее дерево по столбцу sort
func sortTree(trees []Tree) []Tree {
	for _, tree := range trees {
		sort.Slice(tree.Children, func(i, j int) bool {
			if len(tree.Children) > 0 {
				tree.Children = sortTree(tree.Children)
			}
			return tree.Children[i].Sort < tree.Children[j].Sort
		})
	}
	return trees
}

// Добавление дерева в БД
func (tree Tree) AddInDB(year int) (err error) {
	Tx := db.InitTx()
	defer Tx.Rollback(context.Background())
	if err = deleteInDB(year, Tx); err != nil {
		log.Println(err)
		return err
	}
	for i := 0; i < len(tree.Children); i++ {
		err = addInDB(year, tree.Children[i], -1, Tx)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	if err = Tx.Commit(context.Background()); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func deleteInDB(year int, Tx pgx.Tx) error {
	_, err := Tx.Exec(context.Background(), `DELETE From eco_pasport."Tree"  where "Year" = $1;`, year)
	if err != nil {
		return err
	}
	return nil
}
func addInDB(year int, tree Tree, idParent int, Tx pgx.Tx) error {
	if len(tree.IdHeader) != 0 {
		_, err := Tx.Exec(context.Background(), `INSERT INTO eco_pasport."Tree_Header" ("Id_Header","Id_Tree") VALUES ($1,$2);`, tree.IdHeader[0], idParent)
		if err != nil {
			return err
		}
	} else {
		var idParentForDB sql.NullInt32
		if idParent < 0 {
			idParentForDB.Valid = false
		} else {
			idParentForDB.Int32 = int32(idParent)
			idParentForDB.Valid = true
		}
		err := Tx.QueryRow(context.Background(), `INSERT INTO eco_pasport."Tree" ("Name","Lvl","Id_Parent","Year",sort) VALUES ($1,$2,$3,$4,$5) RETURNING "Id";`, trimValue(tree.Name), tree.Lvl, idParentForDB, year, tree.Sort).Scan(&tree.Id)
		if err != nil {
			return err
		}
		for i := 0; i < len(tree.Children); i++ {
			err = addInDB(year, tree.Children[i], tree.Id, Tx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
