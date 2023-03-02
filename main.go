package main

import (
	"eco-pasport-input/db"
	"eco-pasport-input/handler"
	"flag"
	"fmt"
	"gitlab.com/gbh007/gojlog"
)

func main() {
	db.ConnectDB()
	staticDir := flag.String("s", "static", "папка с файлами для раздачи веб сервера")
	webPort := flag.Int("p", 80, "порт веб сервера")
	flag.Parse()
	done := handler.Run(fmt.Sprintf(":%d", *webPort), *staticDir)
	gojlog.Info("Сервер запущен")
	<-done
}
