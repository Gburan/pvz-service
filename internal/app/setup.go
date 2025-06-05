package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	_ "pvz-service/docs/rest"
	pvz_v1 "pvz-service/internal/generated/api/v1/proto"
	"pvz-service/internal/grpc/server"
	add_product2 "pvz-service/internal/handler/add_product"
	close_reception2 "pvz-service/internal/handler/close_reception"
	create_pvz2 "pvz-service/internal/handler/create_pvz"
	delete_product2 "pvz-service/internal/handler/delete_product"
	"pvz-service/internal/handler/dummy_login"
	login_user2 "pvz-service/internal/handler/login_user"
	"pvz-service/internal/handler/middleware"
	pvz_info2 "pvz-service/internal/handler/pvz_info"
	register_user2 "pvz-service/internal/handler/register_user"
	start_reception2 "pvz-service/internal/handler/start_reception"
	nower2 "pvz-service/internal/infrastructure/nower"
	"pvz-service/internal/infrastructure/repository/product"
	"pvz-service/internal/infrastructure/repository/pvz"
	"pvz-service/internal/infrastructure/repository/reception"
	"pvz-service/internal/infrastructure/repository/user"
	"pvz-service/internal/logging"
	"pvz-service/internal/metrics"
	"pvz-service/internal/usecase/add_product"
	"pvz-service/internal/usecase/close_reception"
	"pvz-service/internal/usecase/create_pvz"
	"pvz-service/internal/usecase/delete_product"
	"pvz-service/internal/usecase/list_pvzs"
	"pvz-service/internal/usecase/login_user"
	"pvz-service/internal/usecase/pvz_info"
	"pvz-service/internal/usecase/register_user"
	"pvz-service/internal/usecase/start_reception"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	moderatorRights = []middleware.UserRole{middleware.Moderator}
	employeeRights  = []middleware.UserRole{middleware.Employee}
	sigmaRights     = []middleware.UserRole{middleware.Employee, middleware.Moderator}
)

func (a *App) setup(ctx context.Context) error {
	funcs := []func(context.Context) error{
		a.setupLogger,
		a.setupMetrics,
		a.setupValidator,
		a.setupNewPool,
		a.setupRestServer,
		a.setupGrpcServer,
		a.setupMigrationsDB,
	}

	for _, f := range funcs {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) setupLogger(_ context.Context) error {
	var output io.Writer
	switch a.config.App.Logging.Output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		if err := os.MkdirAll(filepath.Dir(a.config.App.Logging.Output), 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		f, err := os.OpenFile(a.config.App.Logging.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		output = f
	}

	var level slog.Level
	switch a.config.App.Logging.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.Handler(slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: level,
	}))
	handler = logging.NewLoggerImpl(handler)
	slog.SetDefault(slog.New(handler))

	return nil
}

func (a *App) setupMetrics(_ context.Context) error {
	prometheus.MustRegister(
		metrics.RestRequestsTotal,
		metrics.RestResponseDuration,
		metrics.RestEndpointsResponsesTotal,
		metrics.CreatedPVZ,
		metrics.CreatedProducts,
		metrics.CreatedReceptions,
	)
	return nil
}

// @title           PVZ service
// @version         1.0
// @description     This is a service for working with PVZ.

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization
func (a *App) setupRestServer(_ context.Context) error {
	nower := nower2.Nower{}

	repPVZ := pvz.NewRepository(a.Pool, nower)
	repReception := reception.NewRepository(a.Pool, nower)
	repProduct := product.NewRepository(a.Pool, nower)
	repUser := user.NewRepository(a.Pool)

	dummy := dummy_login.New(a.config.App.JWTToken, a.validator)
	loginUserusecase := login_user.NewUsecase(repUser)
	loginer := login_user2.New(a.config.App.JWTToken, loginUserusecase, a.validator)
	registerUserusecase := register_user.NewUsecase(repUser)
	register := register_user2.New(registerUserusecase, a.validator)

	createPVZusecase := create_pvz.NewUsecase(repPVZ)
	creator := create_pvz2.New(createPVZusecase, a.validator)
	getPVZusecase := pvz_info.NewUsecase(repPVZ, repReception, repProduct)
	getter := pvz_info2.New(getPVZusecase, a.validator)

	startReceptionusecase := start_reception.NewUsecase(repPVZ, repReception, repProduct)
	starter := start_reception2.New(startReceptionusecase, a.validator)
	closeReceptionusecase := close_reception.NewUsecase(repPVZ, repReception)
	closer := close_reception2.New(closeReceptionusecase, a.validator)

	addProductusecase := add_product.NewUsecase(repPVZ, repReception, repProduct)
	adder := add_product2.New(addProductusecase, a.validator)
	deleteProductusecase := delete_product.NewUsecase(repPVZ, repReception, repProduct)
	deleter := delete_product2.New(deleteProductusecase, a.validator)

	middlewares := func(mustBeOneOfRole []middleware.UserRole, h http.HandlerFunc) http.Handler {
		handler := h
		if len(mustBeOneOfRole) != 0 {
			handler = middleware.AuthMiddleware(a.config.App.JWTToken, mustBeOneOfRole, handler)
		}
		handler = middleware.LoggerMiddleware(handler)
		handler = middleware.PanicMiddleware(handler)
		return handler
	}

	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	v1 := r.PathPrefix("/api/v1").Subrouter()
	v1.Handle("/dummyLogin", middlewares(nil, dummy.DummyLogin)).Methods("POST")
	v1.Handle("/register", middlewares(nil, register.RegisterUser)).Methods("POST")
	v1.Handle("/login", middlewares(nil, loginer.LoginUser)).Methods("POST")
	v1.Handle("/pvz", middlewares(moderatorRights, creator.CreatePVZ)).Methods("POST")
	v1.Handle("/pvz", middlewares(sigmaRights, getter.GetPVZInfo)).Methods("GET")
	v1.Handle("/pvz/{pvzId}/close_last_reception", middlewares(employeeRights, closer.CloseReception)).Methods("POST")
	v1.Handle("/pvz/{pvzId}/delete_last_product", middlewares(employeeRights, deleter.DeleteProduct)).Methods("POST")
	v1.Handle("/receptions", middlewares(employeeRights, starter.StartReception)).Methods("POST")
	v1.Handle("/products", middlewares(employeeRights, adder.AddProduct)).Methods("POST")

	a.restServer = &http.Server{
		Addr:         a.config.Server.Rest.Address,
		ReadTimeout:  a.config.Server.Rest.Connsettings.ReadTimeout,
		WriteTimeout: a.config.Server.Rest.Connsettings.WriteTimeout,
		IdleTimeout:  a.config.Server.Rest.Connsettings.IdleTimeout,
		Handler:      r,
	}

	return nil
}

func (a *App) setupGrpcServer(_ context.Context) error {
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: a.config.Server.GRPC.ConnSettings.MaxConnIdle,
			MaxConnectionAge:  a.config.Server.GRPC.ConnSettings.MaxConnAge,
		}),
	)

	nower := nower2.Nower{}
	repPVZ := pvz.NewRepository(a.Pool, nower)
	listpvzsUsecase := list_pvzs.NewUsecase(repPVZ)
	pvzServer := server.New(listpvzsUsecase)

	pvz_v1.RegisterPVZServiceServer(grpcServer, pvzServer)
	a.grpcServer = grpcServer

	return nil
}

func (a *App) setupNewPool(ctx context.Context) error {
	cfg, err := pgxpool.ParseConfig(a.config.DB.Conn)
	if err != nil {
		return err
	}

	cfg.MaxConns = a.config.DB.PoolSettings.MaxConns
	cfg.MaxConnIdleTime = a.config.DB.PoolSettings.MaxConnIdleTime
	cfg.MinIdleConns = a.config.DB.PoolSettings.MinIdleConns
	cfg.MaxConnLifetime = a.config.DB.PoolSettings.MaxConnLifetime

	a.Pool, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) setupMigrationsDB(_ context.Context) error {
	dsn := flag.String("dsn", a.config.DB.Conn, "PostgreSQL")

	sql, err := goose.OpenDBWithDriver("postgres", *dsn)
	if err != nil {
		return err
	}

	if err = goose.Up(sql, a.config.DB.MigrationsDir); err != nil {
		return err
	}

	return nil
}

func (a *App) setupValidator(_ context.Context) error {
	a.validator = validator.New()
	register := func(category string, allowed []string) error {
		return a.validator.RegisterValidation(category, func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			for _, allow := range allowed {
				if value == allow {
					return true
				}
			}
			return false
		})
	}

	validations := map[string][]string{
		"oneof_category": a.config.App.Validation.AllowedCategories,
		"oneof_city":     a.config.App.Validation.AllowedCities,
		"oneof_user":     a.config.App.Validation.AllowedUsers,
	}

	for cat, catAllow := range validations {
		err := register(cat, catAllow)
		if err != nil {
			return fmt.Errorf("registering category %s error: %v", cat, err)
		}
	}

	return nil
}
