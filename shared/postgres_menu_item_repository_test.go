package shared_test

import (
	"context"
	"golangsevillabar/shared"
	"golangsevillabar/testhelpers"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PostgresMenuItemRepositoryTestSuite struct {
	suite.Suite
	pgContainer        *testhelpers.PostgresContainer
	menuItemRepository shared.MenuItemRepository
	ctx                context.Context
}

func (suite *PostgresMenuItemRepositoryTestSuite) TestReadItemsShouldFailIfCoultNotRetrieveAllUniqueItems() {
	// Given
	itemsToRead := []int{1, 2, 3, 4}
	// When
	items, err := suite.menuItemRepository.ReadItems(suite.ctx, itemsToRead)
	// Then
	assert.Empty(suite.T(), items)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "requested 4 distinct items, but read from DB 3 distinct items", err.Error())
}

func (suite *PostgresMenuItemRepositoryTestSuite) TestReadItems() {
	// Given
	itemsToRead := []int{1, 2, 3}
	// When
	items, err := suite.menuItemRepository.ReadItems(suite.ctx, itemsToRead)
	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), []shared.OrderedItem{
		{
			MenuItem:    1,
			Description: "blue water",
			Price:       1.0,
		},
		{
			MenuItem:    2,
			Description: "red water",
			Price:       2.0,
		},
		{
			MenuItem:    3,
			Description: "green water",
			Price:       3.0,
		},
	}, items)
}

func (suite *PostgresMenuItemRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.T(), suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	menuItemRepositoryPostgres, err := shared.NewPostgresMenuItemRepository(suite.ctx, suite.pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	suite.menuItemRepository = menuItemRepositoryPostgres
}

func (suite *PostgresMenuItemRepositoryTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func TestPostgresMenuItemRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresMenuItemRepositoryTestSuite))
}
