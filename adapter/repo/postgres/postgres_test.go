package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/golangmigrator"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	dbName  = "postgres"
	dbUser  = "postgres"
	dbPass  = "password"
	timeout = 5 * time.Minute
)

var testDb *TestDB

type TestDB struct {
	instance testcontainers.Container
}

func TestMain(m *testing.M) {
	var err error
	testDb, err = newPgTestDb()
	if err != nil {
		log.Printf("error creating test db container: %v", err)
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}

func newPgTestDb() (*TestDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15.4",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
			return fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s search_path=%s sslmode=disable connect_timeout=5",
				host, dbUser, dbPass, port.Int(), dbName, "public")
		}),
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
		},
		Env: map[string]string{
			"POSTGRES_DB":       dbName,
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPass,
		},
	}

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &TestDB{pgContainer}, nil
}

func getHost(db *TestDB) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return db.instance.Host(ctx)
}

func getPort(db *TestDB) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	p, err := db.instance.MappedPort(ctx, "5432")
	if err != nil {
		return 0, err
	}
	return p.Int(), nil
}

func getIsolatedDatabase(t *testing.T) *sql.DB {
	gm := golangmigrator.New("migrations")

	host, err := getHost(testDb)
	if err != nil {
		t.Fatal(err)
	}

	port, err := getPort(testDb)
	if err != nil {
		t.Fatal(err)
	}

	db := pgtestdb.New(t, pgtestdb.Config{
		Host:       host,
		User:       dbUser,
		Password:   dbPass,
		Port:       strconv.Itoa(port),
		Options:    "sslmode=disable",
		Database:   dbName,
		DriverName: "postgres",
	}, gm)

	return db
}
