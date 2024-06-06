package proxy

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zograf/gobserve/core"
)

func getAll(c echo.Context) error {
	cc := c.(*core.CustomContext)

	infos, err := cc.Sr.GetInfos()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, infos)
}

func getByName(c echo.Context) error {
	name := c.Param("name")
	cc := c.(*core.CustomContext)

	infos, err := cc.Sr.GetInfos()
	if err != nil {
		return err
	}

	value, exists := infos[name]
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Service info with that name was not found",
		})
	}

	return c.JSON(http.StatusOK, value)
}

func proxyPass(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
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
