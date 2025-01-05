package service

import (
	queries_mocks "golangsevillabar/queries/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ReadServiceTestSuite struct {
	suite.Suite
	openTabQueries queries_mocks.OpenTabQueries
	readService    *ReadService
}

func (suite *ReadServiceTestSuite) TestActiveTablesHandlerReturnsErrorIfNotPost() {
	// Given
	rr := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "", nil)
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
	request, err := http.NewRequest(http.MethodPost, "", nil)
	assert.NoError(suite.T(), err)
	suite.openTabQueries.On("ActiveTableNumbers").Return([]int{1, 2})

	// When
	suite.readService.activeTablesHandler(rr, request)

	// Then
	assert.Equal(suite.T(), string("200 OK"), rr.Result().Status)
	bytes, err := io.ReadAll(rr.Result().Body)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":true,\"error\":\"\",\"data\":{\"active_tables\":[1,2]}}", string(bytes))
}

func (suite *ReadServiceTestSuite) SetupTest() {
	suite.openTabQueries = *queries_mocks.NewOpenTabQueries(suite.T())
	suite.readService = CreateReadService(1235, &suite.openTabQueries)
}

func TestReadServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ReadServiceTestSuite))
}
