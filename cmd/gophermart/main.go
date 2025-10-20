package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/bootstrap"
	"github.com/ElfAstAhe/cls-gophermart/internal/utils"
)

func main() {
	// app instance
	app := bootstrap.NewApp()
	defer utils.CloseOnly(app)
	logger := app.Log.GetLogger("main")
	//	defer _utl.CloseOnly(logger.(io.Closer))

	// app initialization
	logger.Info("app initialization")
	if err := app.Init(); err != nil {
		logger.Errorf("app initialization failed [%v]", err)

		os.Exit(1)
	}

	// app run
	logger.Info("app running")
	if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Errorf("app run error [%v]", err)
	}

	app.WG.Wait()

	logger.Info("app shutdown")
}
