package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sferawann/go-bank-api/controller"
)

func NewRouter(e *echo.Echo, allController controller.AllController) {

	api := e.Group("/go-bank-api")

	api.POST("/daftar", allController.Create)
	api.POST("/tabung", allController.Tabung)
	api.POST("/tarik", allController.Tarik)
	api.GET("/saldo/:no_rekening", allController.GetSaldo)

}
