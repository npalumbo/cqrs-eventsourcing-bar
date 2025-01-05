package service

import (
	"encoding/json"
	"fmt"
	"golangsevillabar/queries"
	"golangsevillabar/readservice/model"
	"io"
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
	srv.serveMux.HandleFunc("/tabIdForTable", srv.tabIdForTableNumberHandler)
	srv.serveMux.HandleFunc("/tabForTable", srv.tabForTableNumberHandler)
	srv.serveMux.HandleFunc("/todoListForWaiter", srv.todoListForWaiterHandler)
	srv.serveMux.HandleFunc("/invoiceForTable", srv.invoiceForTableNumberHandler)

	srv.httpServer = &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	srv.httpServer.Handler = srv.serveMux
	srv.openTabQueries = openTabQueries

	return srv
}

func (rs *ReadService) invoiceForTableNumberHandler(w http.ResponseWriter, r *http.Request) {

	var request model.InvoiceForTableRequest
	shouldReturn := readRequest(w, r, &request, &model.QueryResponse[any]{})
	if shouldReturn {
		return
	}

	tabInvoice, err := rs.openTabQueries.InvoiceForTable(request.TableNumber)
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

func (rs *ReadService) tabForTableNumberHandler(w http.ResponseWriter, r *http.Request) {

	var request model.TabForTableRequest
	shouldReturn := readRequest(w, r, &request, &model.QueryResponse[any]{})
	if shouldReturn {
		return
	}

	tabStatus, err := rs.openTabQueries.TabForTable(request.TableNumber)
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

func (rs *ReadService) tabIdForTableNumberHandler(w http.ResponseWriter, r *http.Request) {

	var request model.TabIdForTableRequest
	shouldReturn := readRequest(w, r, &request, &model.QueryResponse[any]{})
	if shouldReturn {
		return
	}

	tabId, err := rs.openTabQueries.TabIdForTable(request.TableNumber)
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

func (rs *ReadService) todoListForWaiterHandler(w http.ResponseWriter, r *http.Request) {

	var request model.TodoListForWaiterRequest
	shouldReturn := readRequest(w, r, &request, &model.QueryResponse[any]{})
	if shouldReturn {
		return
	}

	todoListForWaiter := rs.openTabQueries.TodoListForWaiter(request.Waiter)

	todoListForWaiterResponse := model.QueryResponse[model.TodoListForWaiterResponse]{
		Data:  model.TodoListForWaiterResponse{TabItemsForTable: todoListForWaiter},
		OK:    true,
		Error: "",
	}

	returnJsonOk(w, todoListForWaiterResponse)
}

func (rs *ReadService) activeTablesHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
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

func readRequest[T any](w http.ResponseWriter, r *http.Request, data *T, response *model.QueryResponse[any]) (errored bool) {
	if r.Method != http.MethodPost {
		returnJsonError(w, "Method Not Allowed", http.StatusMethodNotAllowed, response)
		errored = true
		return
	}

	if r.Body == nil {
		returnJsonError(w, "Empty body", http.StatusBadRequest, response)
		errored = true
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		returnJsonError(w, "Error reading request body", http.StatusBadRequest, response)
		errored = true
		return
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		returnJsonError(w, "Invalid JSON request", http.StatusBadRequest, response)
		errored = true
		return
	}
	return
}

func returnJsonError[T any](w http.ResponseWriter, error string, code int, response *model.QueryResponse[T]) {
	h := w.Header()

	h.Del("Content-Length")

	h.Set("Content-Type", "application/json")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

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
