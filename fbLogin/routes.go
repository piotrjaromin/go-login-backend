package fbLogin

import (
        "github.com/labstack/echo"
        "github.com/piotrjaromin/login-template/web"
)

//InitRoutes binds http handlers to paths
func InitRoutes(echoEngine *echo.Echo, controller Controller) {

        echoEngine.OPTIONS("/fb/login", web.OptionsMethodHandler)
        echoEngine.POST("/fb/login", controller.Login)
}
