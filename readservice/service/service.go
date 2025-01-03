package service

import (
	"encoding/json"
	"fmt"
	"golangsevillabar/queries"
	"golangsevillabar/readservice/model"
	"log/slog"
	"net/http"
)

type ReadService struct {
	httpServer     *http.Server
	serveMux       *http.ServeMux
	openTabQueries queries.OpenTabQueries
}

func CreateReadService(port int, openTabQueries queries.OpenTabQueries) *ReadService {
	srv := &ReadService{}

	srv.serveMux = http.NewServeMux()
	srv.serveMux.HandleFunc("/activeTableNumbers", srv.activeTablesHandler)

	srv.httpServer = &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	srv.httpServer.Handler = srv.serveMux
	srv.openTabQueries = openTabQueries

	return srv
}

func (rs *ReadService) activeTablesHandler(w http.ResponseWriter, r *http.Request) {

	activeTableNumbersResponse := model.ActiveTableNumbersResponse{
		ActiveTables: rs.openTabQueries.ActiveTableNumbers(),
		OK:           true,
		Error:        "",
	}

	returnJsonOk(w, activeTableNumbersResponse)
}

func (rs *ReadService) Start() error {
	slog.Info(fmt.Sprintf("Read server listening on%s", rs.httpServer.Addr))

	return rs.httpServer.ListenAndServe()
}

func returnJsonOk(w http.ResponseWriter, response interface{}) {
	h := w.Header()

	h.Del("Content-Length")

	h.Set("Content-Type", "application/json")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "error encoding json, command processed sucesssfully", http.StatusOK)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "error writing json response, command processed sucesssfully", http.StatusOK)
	}
}
