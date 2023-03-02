package db

import (
	//"database/sql"
	"context"
	"eco-pasport-input/configs"
	"github.com/jackc/pgx/v4"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	Pool *pgxpool.Pool
)

func ConnectDB() {
	var err error
	Pool, err = pgxpool.Connect(context.Background(), configs.GetConfig().ConnStr)
	if err != nil {
		log.Println(err)
	}
	if err = Pool.Ping(context.Background()); err != nil {
		log.Println(err)
	}
}
func InitTx() pgx.Tx {
	Tx, err := Pool.Begin(context.Background())
	if err != nil {
		log.Println(err)
		return nil
	}
	return Tx
}
