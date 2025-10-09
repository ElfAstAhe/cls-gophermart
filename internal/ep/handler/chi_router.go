package handler

import (
	"net/http"
	"time"

	_cfg "github.com/ElfAstAhe/cls-gophermart/internal/app/config"
	_log "github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	_cpr "github.com/ElfAstAhe/cls-gophermart/internal/ep/middleware/compress"
	_hlg "github.com/ElfAstAhe/cls-gophermart/internal/ep/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type AppChiRouter struct {
	router *chi.Mux
	log    _log.AppLogger
	conf   *_cfg.Config
}

func NewChiRouter(config *_cfg.Config, logger _log.AppLogger) *AppChiRouter {
	res := &AppChiRouter{
		router: chi.NewRouter(),
		log:    logger.GetLogger("app router"),
		conf:   config,
	}

	res.setupMiddleware(logger.GetLogger("middleware"))
	res.setupRoutes()

	return res
}

func (cr *AppChiRouter) GetRouter() http.Handler {
	return cr.router
}

func (cr *AppChiRouter) setupMiddleware(logger _log.AppLogger) {
	cr.router.Use(middleware.RequestID)
	cr.router.Use(middleware.RealIP)
	cr.router.Use(_cpr.CustomCompress(_cpr.DefaultCompressionLevel, _cpr.ContentTypeApplicationJSON, _cpr.ContentTypeTextHTML))
	cr.router.Use(_cpr.CustomDecompress)
	cr.router.Use(_hlg.NewHTTPLoggerMiddleware(logger).CustomInfoHTTPLogger)
	cr.router.Use(middleware.Recoverer)
	cr.router.Use(middleware.Timeout(60 * time.Second))
}

func (cr *AppChiRouter) setupRoutes() {
	// api
	cr.router.Route("/api", func(r chi.Router) {
		// user
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", cr.postApiUserRegister)     // POST /api/user/register
			r.Post("/login", cr.postApiUserLogin)           // POST /api/user/login
			r.Get("/withdrawals", cr.getApiUserWithdrawals) // GET  /api/user/withdrawals
			// orders
			r.Route("/orders", func(r chi.Router) {
				r.Post("/", cr.postApiUserOrders) // POST /api/user/orders
				r.Get("/", cr.getApiUserOrders)   // GET  /api/user/orders
			})
			// balance
			r.Route("/balance", func(r chi.Router) {
				r.Get("/", cr.getApiUserBalance)                   // GET  /api/user/balance
				r.Post("/withdraw", cr.postApiUserBalanceWithdraw) // POST /api/user/balance/withdraw
			})
		})
	})
}
