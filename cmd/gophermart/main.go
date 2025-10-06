package main

import (
	"io"
	"os"

	_bs "github.com/ElfAstAhe/cls-gophermart.git/internal/app/bootstrap"
	_utl "github.com/ElfAstAhe/cls-gophermart.git/internal/utils"
)

func main() {
	// app instance
	app := _bs.NewApp()
	defer _utl.CloseOnly(app)
	logger := app.Log.GetLogger("main")
	defer _utl.CloseOnly(logger.(io.Closer))

	// app initialization
	logger.Info("app initialization")
	if err := app.Init(); err != nil {
		logger.Errorf("app initialization failed [%v]", err)

		os.Exit(1)
	}

	// app run
	logger.Info("app running")
	if err := app.Run(); err != nil {
		logger.Errorf("app run error [%v]", err)

		os.Exit(1)
	}

	logger.Info("app shutdown")
}
