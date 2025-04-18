package main

import (
	"github.com/labstack/echo/v4"
	"github.com/sferawann/go-bank-api/config"
	"github.com/sferawann/go-bank-api/controller"
	"github.com/sferawann/go-bank-api/repository"
	"github.com/sferawann/go-bank-api/router"
	"github.com/sferawann/go-bank-api/usecase"
	"github.com/sferawann/go-bank-api/utils"
)

func main() {
	utils.SetupLogger()

	config.InitDB()
	db := config.DB

	nasabahRepo := repository.NewNasabahRepository(db)
	rekeningRepo := repository.NewRekeningRepository(db)
	transaksiRepo := repository.NewTransaksiRepository(db)

	usecase := usecase.NewUsecase(nasabahRepo, rekeningRepo, transaksiRepo)
	controller := controller.NewController(usecase)

	e := echo.New()
	router.NewRouter(e, controller)

	utils.Log.Infof("Aplikasi berjalan di port :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
