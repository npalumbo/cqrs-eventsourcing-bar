package testhelpers

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(t testing.TB, ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	testcontainers.CleanupContainer(t, pgContainer)
	require.NoError(t, err)
	if err != nil {
		return nil, err
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	migrationsDir, err := filepath.Abs(filepath.Join("..", "migrations"))
	if err != nil {
		return nil, err
	}

	m, err := migrate.New("file://"+migrationsDir, connStr)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	seedPath := filepath.Join("..", "testdata", "seed.sql")
	if err := applySeedSQL(ctx, connStr, seedPath); err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}

func applySeedSQL(ctx context.Context, connStr string, seedPath string) error {
	seedData, err := os.ReadFile(seedPath)
	if err != nil {
		return err
	}

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, string(seedData))
	return err
}
