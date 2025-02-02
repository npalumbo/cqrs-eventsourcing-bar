package service

import (
	"bytes"
	"context"
	"cqrseventsourcingbar/commands"
	commands_mocks "cqrseventsourcingbar/commands/mocks"
	"cqrseventsourcingbar/shared"
	shared_mocks "cqrseventsourcingbar/shared/mocks"
	"cqrseventsourcingbar/writeservice/model"
	"encoding/json"
	"errors"
	"io"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WriteServiceTestSuite struct {
	suite.Suite
	menuItemRepository *shared_mocks.MenuItemRepository
	commandDispatcher  *commands_mocks.CommandDispatcher
	writeService       *WriteService
	ctx                context.Context
}

func (suite *WriteServiceTestSuite) TestOpenTabHandlerReturnsErrorIfNotPost() {
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
func (suite *WriteServiceTestSuite) TestOpenTabHandlerReturnsErrorIfEmptyBody() {
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

func (suite *WriteServiceTestSuite) TestOpenTabHandlerReturnsErrorIfInvalidJson() {

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

func (suite *WriteServiceTestSuite) TestOpenTabHandlerReturnsErrorIfDispatcherReturnsError() {

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

func (suite *WriteServiceTestSuite) TestOpenTabHandlerReturnsOkIfNoError() {

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

func (suite *WriteServiceTestSuite) TestPlaceOrderHandlerReturnsErrorIfNotPost() {
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
func (suite *WriteServiceTestSuite) TestPlaceOrderHandlerReturnsErrorIfEmptyBody() {
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

func (suite *WriteServiceTestSuite) TestPlaceOrderHandlerReturnsErrorIfInvalidJson() {

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

func (suite *WriteServiceTestSuite) TestPlaceOrderHandlerReturnsErrorIfCannotParseId() {

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

func (suite *WriteServiceTestSuite) TestPlaceOrderHandlerReturnsErrorIfMenuItemRepositoryReturnsError() {

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

	suite.menuItemRepository.On("ReadItems", suite.ctx, []int{1}).Return([]shared.MenuItem{}, errors.New("error from menuItemRepo"))

	// When
	suite.writeService.placeOrderHandler(rr, request)

	// Then
	assert.Equal(suite.T(), "400 Bad Request", rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":false,\"error\":\"could not read items from DB\"}", string(bytes))
}

func (suite *WriteServiceTestSuite) TestPlaceOrderHandlerReturnsErrorIfDispatcherReturnsError() {

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

	suite.menuItemRepository.On("ReadItems", suite.ctx, []int{1}).Return([]shared.MenuItem{{
		ID:          1,
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

func (suite *WriteServiceTestSuite) TestMarkDrinksServedHandlerReturnsErrorIfNotPost() {
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
func (suite *WriteServiceTestSuite) TestMarkDrinksServedHandlerReturnsErrorIfEmptyBody() {
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

func (suite *WriteServiceTestSuite) TestMarkDrinksServedHandlerReturnsErrorIfInvalidJson() {

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

func (suite *WriteServiceTestSuite) TestMarkDrinksServedHandlerReturnsErrorIfCannotParseId() {

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

func (suite *WriteServiceTestSuite) TestMarkDrinksServedHandlerReturnsErrorIfDispatcherReturnsError() {

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

func (suite *WriteServiceTestSuite) TestMarkDrinksServedHandlerReturnsOkIfNoError() {

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

func (suite *WriteServiceTestSuite) TestCloseTabHandlerReturnsErrorIfNotPost() {
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
func (suite *WriteServiceTestSuite) TestCloseTabHandlerReturnsErrorIfEmptyBody() {
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

func (suite *WriteServiceTestSuite) TestCloseTabHandlerReturnsErrorIfInvalidJson() {

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

func (suite *WriteServiceTestSuite) TestCloseTabHandlerReturnsErrorIfCannotParseId() {

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

func (suite *WriteServiceTestSuite) TestCloseTabHandlerReturnsErrorIfDispatcherReturnsError() {

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

func (suite *WriteServiceTestSuite) TestCloseTabHandlerReturnsOkIfNoError() {

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

func (suite *WriteServiceTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.menuItemRepository = shared_mocks.NewMenuItemRepository(suite.T())
	suite.commandDispatcher = commands_mocks.NewCommandDispatcher(suite.T())
	suite.writeService = CreateWriteService(1234, suite.menuItemRepository, suite.commandDispatcher)
}

func TestWriteServiceTestSuite(t *testing.T) {
	suite.Run(t, new(WriteServiceTestSuite))
}
