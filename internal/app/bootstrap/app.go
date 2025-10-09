package bootstrap

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_cfg "github.com/ElfAstAhe/cls-gophermart/internal/app/config"
	_db "github.com/ElfAstAhe/cls-gophermart/internal/app/config/db"
	_log "github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	_rep "github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
	_repi "github.com/ElfAstAhe/cls-gophermart/internal/dal/repository"
	_handler "github.com/ElfAstAhe/cls-gophermart/internal/ep/handler"
	_hnd "github.com/ElfAstAhe/cls-gophermart/internal/ep/handler"
	_migr "github.com/ElfAstAhe/cls-gophermart/migrations"
)

type App struct {
	DB           _db.DB
	Log          _log.AppLogger
	Conf         *_cfg.Config
	Router       _handler.AppRouter
	withdrawRepo _rep.WithdrawRepository
	orderRepo    _rep.OrderRepository
	accountRepo  _rep.AccountRepository
	userRepo     _rep.UserRepository
}

func NewApp() *App {
	return &App{
		Log: _log.NewStartupZapLogger(),
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
	if err := _db.CloseDB(app.DB); err != nil {
		return err
	}

	if err := app.Log.Close(); err != nil {
		return err
	}

	return nil
}

func (app *App) loadConfig() error {
	appConf := _cfg.NewConfig()
	err := appConf.LoadConfig()
	if err != nil {
		return err
	}
	app.Conf = appConf

	return nil
}

func (app *App) initLogger() error {
	fullLogger, err := _log.NewZapLogger(app.Conf.LogLevel, app.Conf.LogFilePath)
	if err != nil {
		return err
	}
	app.Log = fullLogger

	return nil
}

func (app *App) initDatabase() error {
	res, err := _db.NewDB(_cfg.DBKindPostgres, app.Conf.DBDsn)
	if err != nil {
		return err
	}
	app.DB = res

	return nil
}

func (app *App) migrateDatabase() error {
	migrator, err := _migr.NewGooseDBMigrator(context.Background(), app.DB.GetDB(), app.Log.GetLogger("migration"))
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
	app.withdrawRepo = _repi.NewWithdrawPgRepository(app.DB)
	app.orderRepo = _repi.NewOrderPgRepository(app.DB)
	app.accountRepo = _repi.NewAccountPgRepository(app.DB, app.withdrawRepo, app.orderRepo)
	app.userRepo = _repi.NewUserPgRepository(app.DB, app.accountRepo)

	return nil
}

func (app *App) initStartupServices() error {
	// ToDo: implement

	return nil
}

func (app *App) initRouter() error {
	app.Router = _hnd.NewChiRouter(app.Conf, app.Log)

	return nil
}

func (app *App) gracefulShutdown() {
	// channel
	sig := make(chan os.Signal, 1)
	// register channel signals
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// awaiting signal
	<-sig

	//if err := _db.CloseDB(app.DB); err != nil {
	//	app.Log.Errorf("Error closing database: [%v]", err)
	//}

	app.Log.Info("Graceful shutdown server done")

	os.Exit(0)
}
