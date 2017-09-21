package accounts

import (
        "github.com/op/go-logging"
        "github.com/labstack/echo"
        "github.com/piotrjaromin/go-login-backend/web"
        "github.com/piotrjaromin/go-login-backend/security"
)

//InitRoutes binds http handlers to paths
func InitRoutes(echoEngine *echo.Echo, controller Controller, security security.Security) {

        log := logging.MustGetLogger("[initRoutes]")
        log.Debugf("Creating routes")

        accountGroup := echoEngine.Group("/accounts")

        //Accounts endpoints
        accountGroup.OPTIONS("", web.OptionsMethodHandler)
        accountGroup.OPTIONS("/", web.OptionsMethodHandler)
        accountGroup.POST("/", controller.Create)
        accountGroup.POST("", controller.Create)

        accountGroup.OPTIONS("/:id/confirm", web.OptionsMethodHandler)
        accountGroup.GET("/:id/confirm", controller.ConfirmAccount)

        accountGroup.OPTIONS("/:id/reset", web.OptionsMethodHandler)
        accountGroup.POST("/:id/reset", controller.ResetPassword)
        accountGroup.PUT("/:id/reset", controller.ConfirmResetPassword)

        accountGroup.OPTIONS("/:id", web.OptionsMethodHandler)
        accountGroup.Use(security.SecuredById("username", "username", false))
        accountGroup.GET("/:id", controller.GetByID)
}