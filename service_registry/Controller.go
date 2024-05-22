package serviceregistry

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zograf/gobserve/core"
)

func getAll(c echo.Context) error {
	cc := c.(*core.CustomContext)
	return c.JSON(http.StatusOK, cc.Sr.GetInfos())
}

func getByName(c echo.Context) error {
	name := c.Param("name")
	cc := c.(*core.CustomContext)

	value, exists := cc.Sr.GetInfos()[name]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Service info with that name was not found",
		})
	}

	return c.JSON(http.StatusOK, value)
}

func register(c echo.Context) error {
	si := &core.ServiceInfo{}

	if err := c.Bind(si); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	cc := c.(*core.CustomContext)
	if err := cc.Sr.AddServiceInfo(si); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, si)
}
