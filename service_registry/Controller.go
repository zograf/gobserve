package serviceregistry

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func getAll(c echo.Context) error {
	cc := c.(*CustomContext)
	return c.JSON(http.StatusOK, cc.Sr.Infos)
}

func getByName(c echo.Context) error {
	name := c.Param("name")
	cc := c.(*CustomContext)

	value, exists := cc.Sr.Infos[name]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Service info with that name was not found",
		})
	}

	return c.JSON(http.StatusOK, value)
}

func register(c echo.Context) error {
	si := &ServiceInfo{}

	if err := c.Bind(si); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	cc := c.(*CustomContext)
	if err := cc.Sr.AddServiceInfo(si); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, si)
}
