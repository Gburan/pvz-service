package app

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"pvz-service/internal/config"
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
	"pvz-service/internal/infrastructure/repository/product"
	"pvz-service/internal/infrastructure/repository/pvz"
	"pvz-service/internal/infrastructure/repository/reception"
	"pvz-service/internal/infrastructure/repository/user"
	"pvz-service/internal/usecase/add_product"
	"pvz-service/internal/usecase/close_reception"
	"pvz-service/internal/usecase/create_pvz"
	"pvz-service/internal/usecase/delete_product"
	"pvz-service/internal/usecase/login_user"
	"pvz-service/internal/usecase/pvz_info"
	"pvz-service/internal/usecase/register_user"
	"pvz-service/internal/usecase/start_reception"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var (
	moderatorRights = []middleware.UserRole{middleware.Moderator}
	employeeRights  = []middleware.UserRole{middleware.Employee}
	bothRights      = []middleware.UserRole{middleware.Employee, middleware.Moderator}
)

type App struct {
	server    http.Server
	config    config.Config
	validator *validator.Validate
	pool      *pgxpool.Pool
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func (a *App) Stop() error {
	log.Println("Gracefully shutdown...")
	return a.server.Shutdown(context.Background())
}

func NewApp(ctx context.Context, cfg config.Config) (*App, error) {
	a := &App{
		config: cfg,
	}

	err := a.setup(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) setup(ctx context.Context) error {
	funcs := []func(context.Context) error{
		a.setupValidator,
		a.newPool,
		a.setupHttpServer,
		a.runMigrationsDB,
	}

	for _, f := range funcs {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) setupHttpServer(_ context.Context) error {
	dummy := dummy_login.New(a.config.App.JWTToken, a.validator)

	repPVZ := pvz.NewRepository(a.pool)
	repReception := reception.NewRepository(a.pool)
	repProduct := product.NewRepository(a.pool)
	repUser := user.NewRepository(a.pool)

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

	middlewares := func(acceptableRoles []middleware.UserRole, h http.HandlerFunc) http.Handler {
		handler := h
		if len(acceptableRoles) != 0 {
			handler = middleware.AuthMiddleware(a.config.App.JWTToken, acceptableRoles, handler)
		}
		handler = middleware.LoggerMiddleware(handler)
		handler = middleware.PanicMiddleware(handler)
		return handler
	}

	r := mux.NewRouter()
	r.Handle("/dummyLogin", middlewares(nil, dummy.DummyLogin)).Methods("POST")
	r.Handle("/register", middlewares(nil, register.RegisterUser)).Methods("POST")
	r.Handle("/login", middlewares(nil, loginer.LoginUser)).Methods("POST")
	r.Handle("/pvz", middlewares(moderatorRights, creator.CreatePVZ)).Methods("POST")
	r.Handle("/pvz", middlewares(bothRights, getter.GetPVZInfo)).Methods("GET")
	r.Handle("/pvz/{pvzId}/close_last_reception", middlewares(employeeRights, closer.CloseReception)).Methods("POST")
	r.Handle("/pvz/{pvzId}/delete_last_product", middlewares(employeeRights, deleter.DeleteProduct)).Methods("POST")
	r.Handle("/receptions", middlewares(employeeRights, starter.StartReception)).Methods("POST")
	r.Handle("/products", middlewares(employeeRights, adder.AddProduct)).Methods("POST")

	log.Println("starting on:" + a.config.Server.Address)
	a.server = http.Server{
		Addr:    a.config.Server.Address,
		Handler: r,
	}

	return nil
}

func (a *App) newPool(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, a.config.DB.Conn)
	if err != nil {
		return err
	}
	a.pool = pool

	return nil
}

func (a *App) runMigrationsDB(_ context.Context) error {
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
		"oneof_category": a.config.App.AllowedCategories,
		"oneof_city":     a.config.App.AllowedCities,
		"oneof_user":     a.config.App.AllowedUsers,
	}

	for cat, catAllow := range validations {
		err := register(cat, catAllow)
		if err != nil {
			return fmt.Errorf("registering category %s error: %v", cat, err)
		}
	}

	return nil
}
