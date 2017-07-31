package login

import (
        "github.com/labstack/echo"
        "github.com/piotrjaromin/login-template/web"
        "github.com/piotrjaromin/login-template/accounts"
        "github.com/op/go-logging"
)

//Controller struct with login and logout functions
type Controller struct {
        Login  func(c echo.Context) error
        Logout func(c echo.Context) error
}

//Create controller responsible for logging in user
func Create(loginService Service) Controller {

        var log = logging.MustGetLogger("[LoginController]")
        
        login := func(c echo.Context) error {
                username := c.FormValue("username")
                pass := c.FormValue("password")
                
                token, err := loginService.Login(username, accounts.Password(pass))

                if err == ErrNotFoundAccount {
                        return web.NotFoundResponse(c)
                }
                
                if err == ErrBadCredentials {
                        return web.BadRequestResponse(c, "Wrong credentials")
                }

                if err == ErrNotConfirmedAccount {
                        return web.ConflictResponse(c, "Account is not confirmed")
                }
                
                if err != nil {
                        return web.LogAndReturnInternalError(c, "Error while performing login", err)
                }

                log.Infof("Generated token is %s", token.Token)
                return c.JSON(200, token)
        }

        logout := func(c echo.Context) error {
                return c.String(200, "")
        }

        return Controller{
                Login: login,
                Logout: logout,
        }
}