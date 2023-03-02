package handler

import (
	"eco-pasport-input/handler/header"
	"eco-pasport-input/handler/parseExcel"
	"eco-pasport-input/handler/permission"
	"eco-pasport-input/handler/regions"
	"eco-pasport-input/handler/rows"
	"eco-pasport-input/handler/tables"
	"eco-pasport-input/handler/tree"
	"eco-pasport-input/handler/year"
	"net/http"
	"time"

	"gitlab.com/gbh007/gojlog"
)

func Run(addr string, staticDir string) <-chan struct{} {
	gojlog.Info("Заходим в Run")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))

	permission.ApplyMux(mux)
	header.ApplyMux(mux)
	regions.ApplyMux(mux)
	year.ApplyMux(mux)
	tables.ApplyMux(mux)
	tree.ApplyMux(mux)
	rows.ApplyMux(mux)
	parseExcel.ApplyMux(mux)

	// создание объекта сервера
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 3 * time.Minute,
		IdleTimeout:  1 * time.Minute,
	}
	done := make(chan struct{})
	go func() {
		if err := server.ListenAndServe(); err != nil {
			gojlog.Error(err)
		}
		close(done)
	}()
	return done
}
