package login

import (
        "github.com/labstack/echo"
        "github.com/piotrjaromin/login-template/web"
)

//InitRoutes binds http handlers to paths
func InitRoutes(echoEngine *echo.Echo, controller Controller) {

        //Login endpoints
        echoEngine.OPTIONS("/login", web.OptionsMethodHandler)
        echoEngine.OPTIONS("logout", web.OptionsMethodHandler)
        echoEngine.POST("/login", controller.Login)
        echoEngine.POST("/logout", controller.Logout)

}