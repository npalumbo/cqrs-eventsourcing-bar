package service

import (
	"encoding/json"
	"fmt"
	"golangsevillabar/queries"
	"golangsevillabar/readservice/model"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
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
	srv.serveMux.HandleFunc("/tabIdForTable", srv.tabIdForTableNumberHandler)
	srv.serveMux.HandleFunc("/tabForTable", srv.tabForTableNumberHandler)
	srv.serveMux.HandleFunc("/invoiceForTable", srv.invoiceForTableNumberHandler)
	srv.serveMux.HandleFunc("/todoListForWaiter", srv.todoListForWaiterHandler)

	srv.httpServer = &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	srv.httpServer.Handler = srv.serveMux
	srv.openTabQueries = openTabQueries

	return srv
}

func (rs *ReadService) activeTablesHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		returnJsonError(w, "Method Not Allowed", http.StatusMethodNotAllowed, &model.QueryResponse[any]{})
		return
	}

	activeTableNumbersResponse := model.QueryResponse[model.ActiveTableNumbersResponse]{
		Data: model.ActiveTableNumbersResponse{
			ActiveTables: rs.openTabQueries.ActiveTableNumbers(),
		},
		OK:    true,
		Error: "",
	}

	returnJsonOk(w, activeTableNumbersResponse)
}

func (rs *ReadService) tabIdForTableNumberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		returnJsonError(w, "Method Not Allowed", http.StatusMethodNotAllowed, &model.QueryResponse[any]{})
		return
	}

	q := r.URL.Query()

	tableNumber, errored := readTableNumber(q, w)
	if errored {
		return
	}

	tabId, err := rs.openTabQueries.TabIdForTable(tableNumber)
	if err != nil {
		returnJsonError(w, fmt.Sprintf("Error processing tabIdForTable request: %v", err), http.StatusInternalServerError, &model.QueryResponse[any]{})
		return
	}
	tabIdForTableResponse := model.QueryResponse[model.TabIdForTableResponse]{
		Data:  model.TabIdForTableResponse{TabId: tabId.String()},
		OK:    true,
		Error: "",
	}

	returnJsonOk(w, tabIdForTableResponse)
}

func (rs *ReadService) tabForTableNumberHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		returnJsonError(w, "Method Not Allowed", http.StatusMethodNotAllowed, &model.QueryResponse[any]{})
		return
	}

	q := r.URL.Query()

	tableNumber, errored := readTableNumber(q, w)
	if errored {
		return
	}
	tabStatus, err := rs.openTabQueries.TabForTable(tableNumber)
	if err != nil {
		returnJsonError(w, fmt.Sprintf("Error processing tabForTable request: %v", err), http.StatusInternalServerError, &model.QueryResponse[any]{})
		return
	}
	tabForTableResponse := model.QueryResponse[model.TabForTableResponse]{
		Data:  model.TabForTableResponse{TabStatus: tabStatus},
		OK:    true,
		Error: "",
	}

	returnJsonOk(w, tabForTableResponse)
}

func (rs *ReadService) invoiceForTableNumberHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		returnJsonError(w, "Method Not Allowed", http.StatusMethodNotAllowed, &model.QueryResponse[any]{})
		return
	}

	q := r.URL.Query()

	tableNumber, errored := readTableNumber(q, w)
	if errored {
		return
	}

	tabInvoice, err := rs.openTabQueries.InvoiceForTable(tableNumber)
	if err != nil {
		returnJsonError(w, fmt.Sprintf("Error processing invoiceForTable request: %v", err), http.StatusInternalServerError, &model.QueryResponse[any]{})
		return
	}
	invoiceForTableResponse := model.QueryResponse[model.InvoiceForTableResponse]{
		Data:  model.InvoiceForTableResponse{TabInvoice: tabInvoice},
		OK:    true,
		Error: "",
	}

	returnJsonOk(w, invoiceForTableResponse)
}

func (rs *ReadService) todoListForWaiterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		returnJsonError(w, "Method Not Allowed", http.StatusMethodNotAllowed, &model.QueryResponse[any]{})
		return
	}

	waiter := r.URL.Query().Get("waiter")

	todoListForWaiter := rs.openTabQueries.TodoListForWaiter(waiter)

	todoListForWaiterResponse := model.QueryResponse[model.TodoListForWaiterResponse]{
		Data:  model.TodoListForWaiterResponse{TabItemsForTable: todoListForWaiter},
		OK:    true,
		Error: "",
	}

	returnJsonOk(w, todoListForWaiterResponse)
}

func readTableNumber(q url.Values, w http.ResponseWriter) (int, bool) {
	tableNumberStr := q.Get("table_number")

	if tableNumberStr == "" {
		returnJsonError(w, "table_number is required", http.StatusBadRequest, &model.QueryResponse[any]{})
		return 0, true
	}

	tableNumber, err := strconv.ParseInt(tableNumberStr, 10, 64)
	if err != nil {
		returnJsonError(w, fmt.Sprintf("Error reading table_number: %v", err), http.StatusBadRequest, &model.QueryResponse[any]{})
		return 0, true
	}
	return int(tableNumber), false
}

func (rs *ReadService) Start() error {
	slog.Info(fmt.Sprintf("Read server listening on%s", rs.httpServer.Addr))

	return rs.httpServer.ListenAndServe()
}

func returnJsonOk(w http.ResponseWriter, response interface{}) {
	setHeaders(w, http.StatusOK)

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

func returnJsonError[T any](w http.ResponseWriter, error string, code int, response *model.QueryResponse[T]) {
	setHeaders(w, code)

	response.Error = error
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("error encoding json, original error: %s", error), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("error writing json response, original error: %s", error), http.StatusInternalServerError)
	}
}

func setHeaders(w http.ResponseWriter, code int) {
	h := w.Header()

	h.Del("Content-Length")

	h.Set("Content-Type", "application/json")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
}
