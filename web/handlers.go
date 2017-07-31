package web

import "github.com/labstack/echo"

func OptionsMethodHandler(c echo.Context) error {
        return c.String(200, "")
}
