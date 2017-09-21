package fbLogin

import (
	"github.com/labstack/echo"
	"github.com/op/go-logging"
	"github.com/piotrjaromin/go-login-backend/web"
)

//Controller for facebook login endpoint
type Controller struct {
	Login func(c echo.Context) error
}

//Create with fb login handler
func Create(service Service) Controller {

	var log = logging.MustGetLogger("[FbLoginController]")

	//Login with use of facebook
	login := func(c echo.Context) error {

		log.Debug("Starting login with facebook")
		fbToken := Token{}
		c.Bind(&fbToken)

		errors := fbToken.Validate()
		if len(errors) > 0 {
			return web.BadRequestResponseWithDetails(c, "Invalid payload", errors)
		}

		appToken, err := service.Login(fbToken)
		if err != nil {
			return web.LogAndReturnInternalError(c, "Could not obtain facebook token", err)
		}

		log.Debugf("Generated token from fbController is %+v", appToken)
		return c.JSON(200, appToken)
	}

	return Controller{
		Login: login,
	}
}
