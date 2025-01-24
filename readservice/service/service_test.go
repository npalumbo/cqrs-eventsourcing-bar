package service

import (
	"errors"
	"golangsevillabar/queries"
	queries_mocks "golangsevillabar/queries/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ReadServiceTestSuite struct {
	suite.Suite
	openTabQueries queries_mocks.OpenTabQueries
	readService    *ReadService
}

func (suite *ReadServiceTestSuite) TestActiveTablesHandlerReturnsErrorIfNotGet() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.readService.activeTablesHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("405 Method Not Allowed"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Method Not Allowed\",\"data\":null}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestActiveTablesHandler() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(suite.T(), err)
	suite.openTabQueries.On("ActiveTableNumbers").Return([]int{1, 2})

	// When
	suite.readService.activeTablesHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("200 OK"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":true,\"error\":\"\",\"data\":[1,2]}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestTabIdForTableNumberReturnsErrorIfNotGet() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.readService.tabIdForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("405 Method Not Allowed"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Method Not Allowed\",\"data\":null}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestTabIdForTableNumberReturnsErrorIfNoTableNumber() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.readService.tabIdForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("400 Bad Request"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"table_number is required\",\"data\":null}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestTabIdForTableNumberReturnsErrorIfBadTableNumber() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "?table_number=NOT_A_NUMBER", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.readService.tabIdForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("400 Bad Request"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Error reading table_number: strconv.ParseInt: parsing \\\"NOT_A_NUMBER\\\": invalid syntax\",\"data\":null}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestTabIdForTableNumberErrorIfQueryErrors() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "?table_number=19", nil)
	assert.NoError(suite.T(), err)
	suite.openTabQueries.On("TabIdForTable", 19).Return(nil, errors.New("fake error"))

	// When
	suite.readService.tabIdForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("500 Internal Server Error"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Error processing tabIdForTable request: fake error\",\"data\":null}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestTabIdForTableNumber() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "?table_number=19", nil)
	assert.NoError(suite.T(), err)
	tabId, _ := ksuid.Parse("2qPTBJCN6ib7iJ6WaIVvoSmySSV")
	suite.openTabQueries.On("TabIdForTable", 19).Return(tabId, nil)

	// When
	suite.readService.tabIdForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("200 OK"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":true,\"error\":\"\",\"data\":\"2qPTBJCN6ib7iJ6WaIVvoSmySSV\"}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestTabForTableNumberReturnsErrorIfNotGet() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.readService.tabForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("405 Method Not Allowed"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Method Not Allowed\",\"data\":null}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestTabForTableNumberReturnsErrorIfNoTableNumber() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.readService.tabForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("400 Bad Request"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"table_number is required\",\"data\":null}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestTabForTableNumberReturnsErrorIfBadTableNumber() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "?table_number=NOT_A_NUMBER", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.readService.tabForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("400 Bad Request"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Error reading table_number: strconv.ParseInt: parsing \\\"NOT_A_NUMBER\\\": invalid syntax\",\"data\":null}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestTabForTableNumberErrorIfQueryErrors() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "?table_number=19", nil)
	assert.NoError(suite.T(), err)
	suite.openTabQueries.On("TabForTable", 19).Return(queries.TabStatus{}, errors.New("fake error"))

	// When
	suite.readService.tabForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("500 Internal Server Error"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Error processing tabForTable request: fake error\",\"data\":null}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestTabForTableNumber() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "?table_number=19", nil)
	assert.NoError(suite.T(), err)
	suite.openTabQueries.On("TabForTable", 19).Return(queries.TabStatus{
		TabID:       "2qPTBJCN6ib7iJ6WaIVvoSmySSV",
		TableNumber: 19,
		ToServe:     []queries.TabItem{},
		Served:      []queries.TabItem{},
	}, nil)

	// When
	suite.readService.tabForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("200 OK"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":true,\"error\":\"\",\"data\":{\"tab_id\":\"2qPTBJCN6ib7iJ6WaIVvoSmySSV\",\"table_number\":19,\"to_serve\":[],\"served\":[]}}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestInvoiceForTableReturnsErrorIfNotGet() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.readService.invoiceForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("405 Method Not Allowed"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Method Not Allowed\",\"data\":null}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestInvoiceForTableReturnsErrorIfNoTableNumber() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.readService.invoiceForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("400 Bad Request"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"table_number is required\",\"data\":null}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestInvoiceForTableReturnsErrorIfBadTableNumber() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "?table_number=NOT_A_NUMBER", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.readService.invoiceForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("400 Bad Request"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Error reading table_number: strconv.ParseInt: parsing \\\"NOT_A_NUMBER\\\": invalid syntax\",\"data\":null}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestInvoiceForTableErrorIfQueryErrors() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "?table_number=19", nil)
	assert.NoError(suite.T(), err)
	suite.openTabQueries.On("InvoiceForTable", 19).Return(queries.TabInvoice{}, errors.New("fake error"))

	// When
	suite.readService.invoiceForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("500 Internal Server Error"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Error processing invoiceForTable request: fake error\",\"data\":null}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestInvoiceForTable() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "?table_number=19", nil)
	assert.NoError(suite.T(), err)
	suite.openTabQueries.On("InvoiceForTable", 19).Return(queries.TabInvoice{
		TabID:            "2qPTBJCN6ib7iJ6WaIVvoSmySSV",
		TableNumber:      19,
		Items:            []queries.TabItem{},
		Total:            0,
		HasUnservedItems: false,
	}, nil)

	// When
	suite.readService.invoiceForTableNumberHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("200 OK"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":true,\"error\":\"\",\"data\":{\"tab_id\":\"2qPTBJCN6ib7iJ6WaIVvoSmySSV\",\"table_number\":19,\"items\":[],\"total\":0,\"has_unserved_items\":false}}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestTodoListForWaiterReturnsErrorIfNotGet() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.readService.todoListForWaiterHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("405 Method Not Allowed"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Method Not Allowed\",\"data\":null}", string(bytes))
}

func (suite *ReadServiceTestSuite) TestTodoListForWaiter() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "?waiter=w1", nil)
	assert.NoError(suite.T(), err)
	tabItems := make(map[int][]queries.TabItem)
	tabItems[19] = []queries.TabItem{{
		MenuNumber:  1,
		Description: "Blue Water",
		Price:       1.0,
	}}
	suite.openTabQueries.On("TodoListForWaiter", "w1").Return(tabItems, nil)

	// When
	suite.readService.todoListForWaiterHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("200 OK"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":true,\"error\":\"\",\"data\":{\"19\":[{\"menu_number\":1,\"description\":\"Blue Water\",\"price\":1}]}}", string(bytes))
}

func (suite *ReadServiceTestSuite) SetupTest() {
	suite.openTabQueries = *queries_mocks.NewOpenTabQueries(suite.T())
	suite.readService = CreateReadService(1235, &suite.openTabQueries)
}

func TestReadServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ReadServiceTestSuite))
}
