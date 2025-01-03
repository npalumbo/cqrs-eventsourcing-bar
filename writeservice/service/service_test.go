package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"golangsevillabar/commands"
	commands_mocks "golangsevillabar/commands/mocks"
	"golangsevillabar/shared"
	shared_mocks "golangsevillabar/shared/mocks"
	"golangsevillabar/writeservice/model"
	"io"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WriteServieTestSuite struct {
	suite.Suite
	menuItemRepository *shared_mocks.MenuItemRepository
	commandDispatcher  *commands_mocks.CommandDispatcher
	writeService       *WriteService
	ctx                context.Context
}

func (suite *WriteServieTestSuite) TestOpenTabHandlerReturnsErrorIfNotPost() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.openTabHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("405 Method Not Allowed"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Method Not Allowed\"}", string(bytes))
}
func (suite *WriteServieTestSuite) TestOpenTabHandlerReturnsErrorIfEmptyBody() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.openTabHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "400 Bad Request", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Empty body\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestOpenTabHandlerReturnsErrorIfInvalidJson() {

	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader([]byte{1, 2, 3}))
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.openTabHandler(rr, request)

	// Then
	assert.Equal(suite.T(), ("400 Bad Request"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Invalid JSON request\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestOpenTabHandlerReturnsErrorIfDispatcherReturnsError() {

	// Given
	openTabRequest := model.OpenTabRequest{
		TableNumber: 0,
		Waiter:      "",
	}
	json, err := json.Marshal(openTabRequest)
	assert.NoError(suite.T(), err)
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(json))
	assert.NoError(suite.T(), err)

	suite.commandDispatcher.On("DispatchCommand", suite.ctx, mock.Anything).Return(errors.New("error dispatching command"))

	// When
	suite.writeService.openTabHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "500 Internal Server Error", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Error processing openTab request: error dispatching command\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestOpenTabHandlerReturnsOkIfNoError() {

	// Given
	openTabRequest := model.OpenTabRequest{
		TableNumber: 1,
		Waiter:      "w1",
	}
	json, err := json.Marshal(openTabRequest)
	assert.NoError(suite.T(), err)
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(json))
	assert.NoError(suite.T(), err)

	var capturedCommand commands.OpenTab
	suite.commandDispatcher.On("DispatchCommand", suite.ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		capturedCommand = args.Get(1).(commands.OpenTab)
	})

	// When
	suite.writeService.openTabHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "200 OK", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":true,\"error\":\"\"}", string(bytes))
	assert.Equal(suite.T(), 1, capturedCommand.TableNumber)
	assert.Equal(suite.T(), "w1", capturedCommand.Waiter)
	assert.NotNil(suite.T(), capturedCommand.ID)
}

func (suite *WriteServieTestSuite) TestPlaceOrderHandlerReturnsErrorIfNotPost() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.placeOrderHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "405 Method Not Allowed", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Method Not Allowed\"}", string(bytes))
}
func (suite *WriteServieTestSuite) TestPlaceOrderHandlerReturnsErrorIfEmptyBody() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.placeOrderHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "400 Bad Request", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Empty body\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestPlaceOrderHandlerReturnsErrorIfInvalidJson() {

	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader([]byte{1, 2, 3}))
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.placeOrderHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "400 Bad Request", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Invalid JSON request\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestPlaceOrderHandlerReturnsErrorIfCannotParseId() {

	// Given
	placeOrderRequest := model.PlaceOrderRequest{
		TabId:     "?",
		MenuItems: []int{1},
	}
	json, err := json.Marshal(placeOrderRequest)
	assert.NoError(suite.T(), err)
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(json))
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.placeOrderHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "400 Bad Request", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"could not parse id\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestPlaceOrderHandlerReturnsErrorIfMenuItemRepositoryReturnsError() {

	// Given
	placeOrderRequest := model.PlaceOrderRequest{
		TabId:     "2qPTBJCN6ib7iJ6WaIVvoSmySSV",
		MenuItems: []int{1},
	}
	json, err := json.Marshal(placeOrderRequest)
	assert.NoError(suite.T(), err)
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(json))
	assert.NoError(suite.T(), err)

	suite.menuItemRepository.On("ReadItems", suite.ctx, []int{1}).Return([]shared.OrderedItem{}, errors.New("error from menuItemRepo"))

	// When
	suite.writeService.placeOrderHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "400 Bad Request", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"could not read items from DB\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestPlaceOrderHandlerReturnsErrorIfDispatcherReturnsError() {

	// Given
	placeOrderRequest := model.PlaceOrderRequest{
		TabId:     "2qPTBJCN6ib7iJ6WaIVvoSmySSV",
		MenuItems: []int{1},
	}
	json, err := json.Marshal(placeOrderRequest)
	assert.NoError(suite.T(), err)
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(json))
	assert.NoError(suite.T(), err)

	suite.menuItemRepository.On("ReadItems", suite.ctx, []int{1}).Return([]shared.OrderedItem{{
		MenuItem:    1,
		Description: "Blue water",
		Price:       1.0,
	}}, nil)
	suite.commandDispatcher.On("DispatchCommand", suite.ctx, mock.Anything).Return(errors.New("error dispatching command"))

	// When
	suite.writeService.placeOrderHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "500 Internal Server Error", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Error processing placeOrder request: error dispatching command\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestMarkDrinksServedHandlerReturnsErrorIfNotPost() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.markDrinksServedHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "405 Method Not Allowed", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Method Not Allowed\"}", string(bytes))
}
func (suite *WriteServieTestSuite) TestMarkDrinksServedHandlerReturnsErrorIfEmptyBody() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.markDrinksServedHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "400 Bad Request", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Empty body\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestMarkDrinksServedHandlerReturnsErrorIfInvalidJson() {

	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader([]byte{1, 2, 3}))
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.markDrinksServedHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "400 Bad Request", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Invalid JSON request\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestMarkDrinksServedHandlerReturnsErrorIfCannotParseId() {

	// Given
	markDrinksServedRequest := model.MarkDrinksServedRequest{
		TabId:       "?",
		MenuNumbers: []int{1},
	}
	json, err := json.Marshal(markDrinksServedRequest)
	assert.NoError(suite.T(), err)
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(json))
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.markDrinksServedHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "400 Bad Request", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"could not parse id\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestMarkDrinksServedHandlerReturnsErrorIfDispatcherReturnsError() {

	// Given
	markDrinksServed := model.MarkDrinksServedRequest{
		TabId:       "2qPTBJCN6ib7iJ6WaIVvoSmySSV",
		MenuNumbers: []int{1},
	}
	json, err := json.Marshal(markDrinksServed)
	assert.NoError(suite.T(), err)
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(json))
	assert.NoError(suite.T(), err)

	suite.commandDispatcher.On("DispatchCommand", suite.ctx, mock.Anything).Return(errors.New("error dispatching command"))

	// When
	suite.writeService.markDrinksServedHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "500 Internal Server Error", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Error processing markDrinksServed request: error dispatching command\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestMarkDrinksServedHandlerReturnsOkIfNoError() {

	// Given
	markDrinksServedRequest := model.MarkDrinksServedRequest{
		TabId:       "2qPTBJCN6ib7iJ6WaIVvoSmySSV",
		MenuNumbers: []int{1},
	}
	json, err := json.Marshal(markDrinksServedRequest)
	assert.NoError(suite.T(), err)
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(json))
	assert.NoError(suite.T(), err)

	var capturedCommand commands.MarkDrinksServed
	suite.commandDispatcher.On("DispatchCommand", suite.ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		capturedCommand = args.Get(1).(commands.MarkDrinksServed)
	})

	// When
	suite.writeService.markDrinksServedHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "200 OK", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":true,\"error\":\"\"}", string(bytes))
	assert.Equal(suite.T(), []int{1}, capturedCommand.MenuNumbers)
	assert.Equal(suite.T(), "2qPTBJCN6ib7iJ6WaIVvoSmySSV", capturedCommand.ID.String())
}

func (suite *WriteServieTestSuite) TestCloseTabHandlerReturnsErrorIfNotPost() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.closeTabHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "405 Method Not Allowed", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Method Not Allowed\"}", string(bytes))
}
func (suite *WriteServieTestSuite) TestCloseTabHandlerReturnsErrorIfEmptyBody() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", nil)
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.closeTabHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "400 Bad Request", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Empty body\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestCloseTabHandlerReturnsErrorIfInvalidJson() {

	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader([]byte{1, 2, 3}))
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.closeTabHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "400 Bad Request", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Invalid JSON request\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestCloseTabHandlerReturnsErrorIfCannotParseId() {

	// Given
	closeTabRequest := model.CloseTabRequest{
		TabId:      "?",
		AmountPaid: 0.0,
	}
	json, err := json.Marshal(closeTabRequest)
	assert.NoError(suite.T(), err)
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(json))
	assert.NoError(suite.T(), err)

	// When
	suite.writeService.closeTabHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "400 Bad Request", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"could not parse id\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestCloseTabHandlerReturnsErrorIfDispatcherReturnsError() {

	// Given
	closeTabRequest := model.CloseTabRequest{
		TabId:      "2qPTBJCN6ib7iJ6WaIVvoSmySSV",
		AmountPaid: 1.0,
	}
	json, err := json.Marshal(closeTabRequest)
	assert.NoError(suite.T(), err)
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(json))
	assert.NoError(suite.T(), err)

	suite.commandDispatcher.On("DispatchCommand", suite.ctx, mock.Anything).Return(errors.New("error dispatching command"))

	// When
	suite.writeService.closeTabHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "500 Internal Server Error", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"Error processing closeTab request: error dispatching command\"}", string(bytes))
}

func (suite *WriteServieTestSuite) TestCloseTabHandlerReturnsOkIfNoError() {

	// Given
	closeTabRequest := model.CloseTabRequest{
		TabId:      "2qPTBJCN6ib7iJ6WaIVvoSmySSV",
		AmountPaid: 1.0,
	}
	json, err := json.Marshal(closeTabRequest)
	assert.NoError(suite.T(), err)
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(json))
	assert.NoError(suite.T(), err)

	var capturedCommand commands.CloseTab
	suite.commandDispatcher.On("DispatchCommand", suite.ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		capturedCommand = args.Get(1).(commands.CloseTab)
	})

	// When
	suite.writeService.closeTabHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "200 OK", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":true,\"error\":\"\"}", string(bytes))
	assert.Equal(suite.T(), 1.0, capturedCommand.AmountPaid)
	assert.Equal(suite.T(), "2qPTBJCN6ib7iJ6WaIVvoSmySSV", capturedCommand.ID.String())
}

func (suite *WriteServieTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.menuItemRepository = shared_mocks.NewMenuItemRepository(suite.T())
	suite.commandDispatcher = commands_mocks.NewCommandDispatcher(suite.T())
	suite.writeService = CreateWriteService(1234, suite.menuItemRepository, suite.commandDispatcher)
}

func TestWriteServieTestSuite(t *testing.T) {
	suite.Run(t, new(WriteServieTestSuite))
}
