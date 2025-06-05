package integrational

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"pvz-service/internal/app"
	"pvz-service/internal/config"

	"github.com/Masterminds/squirrel"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
)

const (
	pgImage = "postgres:16-alpine"

	migrationsPath = "../../migrations"
	cfgPath        = "../../config/config.yaml"
)

var (
	excludeTables = []string{"goose_db_version"}
)

type AppTestSuite struct {
	suite.Suite
	tables    []string
	cfg       config.Config
	app       *app.App
	container testcontainers.Container
}

func (s *AppTestSuite) SetupSuite() {
	cfg := config.MustLoad(cfgPath)
	cfg.DB.MigrationsDir = migrationsPath
	s.cfg = cfg

	ctx := context.Background()
	err := s.createContainer(ctx)
	s.Assert().NoError(err)

	s.app, err = app.NewApp(ctx, s.cfg)
	s.Assert().NoError(err)

	err = s.getTables(ctx)
	s.Assert().NoError(err)

	go func() {
		if err = s.app.RestRun(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.T().Errorf("rest error: %v", err)
		}
	}()
	go func() {
		if err = s.app.GrpcRun(); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			s.T().Errorf("grpc error: %v", err)
		}
	}()

	time.Sleep(2 * time.Second)
}

func (s *AppTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.app.Stop(ctx); err != nil {
		s.T().Errorf("error stop app: %v", err)
	}
	if err := s.container.Terminate(context.Background()); err != nil {
		s.T().Errorf("error terminate container: %v", err)
	}
}

func (s *AppTestSuite) SetupTest() {
	ctx := context.Background()
	truncateSQL := fmt.Sprintf("%s %s %s", "TRUNCATE TABLE", strings.Join(s.tables, ", "), "CASCADE;")
	_, err := s.app.Pool.Exec(ctx, truncateSQL)
	s.Assert().NoError(err)
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

func (s *AppTestSuite) getTables(ctx context.Context) error {
	selectBuilder := squirrel.
		Select("tablename").
		PlaceholderFormat(squirrel.Dollar).
		From("pg_tables").
		Where(squirrel.Eq{"schemaname": "public"})

	if len(excludeTables) > 0 {
		args := make([]interface{}, len(excludeTables))
		for i, v := range excludeTables {
			args[i] = v
		}

		placeholders := strings.Repeat("?,", len(excludeTables))
		placeholders = placeholders[:len(placeholders)-1]

		selectBuilder = selectBuilder.Where(
			fmt.Sprintf("tablename NOT IN (%s)", placeholders),
			args...,
		)
	}

	sql, args, err := selectBuilder.ToSql()
	if err != nil {
		return err
	}

	rows, err := s.app.Pool.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err = rows.Scan(&tableName); err != nil {
			return err
		}
		tables = append(tables, tableName)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	s.tables = tables
	return nil
}

func (s *AppTestSuite) createContainer(ctx context.Context) error {
	u, err := url.Parse(s.cfg.DB.Conn)
	if err != nil {
		return err
	}
	password, _ := u.User.Password()
	if err != nil {
		return err
	}
	_, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return err
	}

	var env = map[string]string{
		"POSTGRES_PASSWORD": u.User.Username(),
		"POSTGRES_USER":     password,
		"POSTGRES_DB":       strings.TrimPrefix(u.Path, "/"),
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        pgImage,
			ExposedPorts: []string{fmt.Sprintf("%s/tcp", port)},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
			Name:         "test-postgres",
		},
		Started: true,
		Reuse:   true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return fmt.Errorf("failed to get container external port: %v", err)
	}

	slog.Info("postgres container ready and running at port: ", "port", mappedPort.Port())

	hostname := u.Hostname()
	u.Host = fmt.Sprintf("%s:%s", hostname, mappedPort.Port())
	s.cfg.DB.Conn = u.String()

	s.container = container

	time.Sleep(time.Second)
	return nil
}
