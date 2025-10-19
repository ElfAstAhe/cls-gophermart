package bootstrap

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/config"
	"github.com/ElfAstAhe/cls-gophermart/internal/app/config/db"
	"github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	irepo "github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/service"
	"github.com/ElfAstAhe/cls-gophermart/internal/dal/repository"
	"github.com/ElfAstAhe/cls-gophermart/internal/ep/facade"
	"github.com/ElfAstAhe/cls-gophermart/internal/ep/handler"
	"github.com/ElfAstAhe/cls-gophermart/internal/utils"
	"github.com/ElfAstAhe/cls-gophermart/migrations"
)

type App struct {
	DB               db.DB
	Log              logger.Logger
	Conf             *config.Config
	withdrawRepo     irepo.WithdrawRepository
	orderRepo        irepo.OrderRepository
	accountRepo      irepo.AccountRepository
	userRepo         irepo.UserRepository
	accountService   service.AccountService
	authService      service.AuthService
	orderPollService service.OrdersPollingService
	usersFacade      facade.UsersFacade
	authFacade       facade.AuthFacade
	Router           handler.AppRouter
}

func NewApp() *App {
	return &App{
		Log: logger.NewStartupZapLogger(),
	}
}

func (app *App) Init() error {
	logger := app.Log.GetLogger("app")
	//    defer _utl.CloseOnly(logger.(io.Closer))

	logger.Info("loading config")
	if err := app.loadConfig(); err != nil {
		return err
	}

	logger.Info("initializing logger")
	if err := app.initLogger(); err != nil {
		return err
	}

	logger.Info("initializing database")
	if err := app.initDatabase(); err != nil {
		return err
	}

	logger.Info("migrate database")
	if err := app.migrateDatabase(); err != nil {
		return err
	}

	logger.Info("initializing dependencies")
	if err := app.initDependencies(); err != nil {
		return err
	}

	logger.Info("initializing startup services")
	if err := app.initStartupServices(); err != nil {
		return err
	}

	logger.Info("initializing http server")
	if err := app.initRouter(); err != nil {
		return err
	}

	return nil
}

func (app *App) Run() error {
	logger := app.Log.GetLogger("app run")
	//    defer _utl.CloseOnly(logger.(io.Closer))
	logger.Info("Starting graceful shutdown go routine...")
	go app.gracefulShutdown()

	logger.Info("Starting server...")
	if err := http.ListenAndServe(app.Conf.HTTP.GetListenerAddr(), app.Router.GetRouter()); err != nil {
		logger.Errorf("Error starting server with error [%v]", err)

		os.Exit(1)
	}

	return nil
}

func (app *App) Close() error {
	if err := app.orderPollService.Stop(); err != nil {
		return err
	}

	if err := db.CloseDB(app.DB); err != nil {
		return err
	}

	if err := app.Log.Close(); err != nil {
		return err
	}

	return nil
}

func (app *App) loadConfig() error {
	appConf := config.NewConfig()
	err := appConf.LoadConfig()
	if err != nil {
		return err
	}
	app.Conf = appConf

	return nil
}

func (app *App) initLogger() error {
	fullLogger, err := logger.NewZapLogger(app.Conf.LogLevel, app.Conf.LogFilePath)
	if err != nil {
		return err
	}
	app.Log = fullLogger

	return nil
}

func (app *App) initDatabase() error {
	res, err := db.NewDB(config.DBKindPostgres, app.Conf.DBDsn)
	if err != nil {
		return err
	}
	app.DB = res

	return nil
}

func (app *App) migrateDatabase() error {
	migrator, err := migrations.NewGooseDBMigrator(context.Background(), app.DB.GetDB(), app.Log.GetLogger("migration"))
	if err != nil {
		return err
	}

	if err := migrator.Initialize(); err != nil {
		return err
	}

	if err := migrator.Up(); err != nil {
		return err
	}

	return nil
}

func (app *App) initDependencies() error {
	// repositories
	app.withdrawRepo = repository.NewWithdrawPgRepository(app.DB)
	app.orderRepo = repository.NewOrderPgRepository(app.DB)
	app.accountRepo = repository.NewAccountPgRepository(app.DB, app.withdrawRepo, app.orderRepo)
	app.userRepo = repository.NewUserPgRepository(app.DB, app.accountRepo)

	// services
	app.orderPollService = service.NewOrdersPollingService(context.Background(), app.Conf.AccrualBaseURI, app.orderRepo, app.Log)
	app.accountService = service.NewAccountServiceImpl(app.orderPollService, app.accountRepo, app.withdrawRepo, app.orderRepo)
	app.authService = service.NewAuthService(app.userRepo)

	// facade
	app.usersFacade = facade.NewUsersFacadeImpl(app.accountService, app.Log)
	app.authFacade = facade.NewAuthFacadeImpl(app.authService, app.Log)

	return nil
}

func (app *App) initStartupServices() error {
	// polling service
	if err := app.orderPollService.Start(); err != nil {
		return err
	}

	return nil
}

func (app *App) initRouter() error {
	app.Router = handler.NewChiRouter(app.Conf, app.usersFacade, app.authFacade, app.Log)

	return nil
}

func (app *App) gracefulShutdown() {
	// channel
	sig := make(chan os.Signal, 1)
	// register channel signals
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// awaiting signal
	<-sig

	utils.CloseOnly(app)

	app.Log.Info("Graceful shutdown server done")

	os.Exit(0)
}
