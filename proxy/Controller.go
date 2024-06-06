package proxy

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func getAll(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func getByName(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func proxyPass(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
