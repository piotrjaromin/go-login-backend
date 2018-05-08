package accounts

import (
	"github.com/labstack/echo"
	"github.com/op/go-logging"
	"github.com/piotrjaromin/go-login-backend/web"
	"net/http"
	"net/url"
)

type Controller struct {
	Create               func(c echo.Context) error
	GetByID              func(c echo.Context) error
	ResetPassword        func(c echo.Context) error
	ConfirmAccount       func(c echo.Context) error
	ConfirmResetPassword func(c echo.Context) error
	Update               func(c echo.Context) error
}

func Create(service Service) Controller {

	var log = logging.MustGetLogger("[AccountController]")

	create := func(c echo.Context) error {

		log.Debug("in create account method")
		account := new(Account)
		if err := c.Bind(account); err != nil {
			return web.BadRequestResponse(c, "Invalid payload")
		}

		secAccount := new(SecuredAccount)
		secAccount.Account = *account

		validationErrors := secAccount.validate()

		if len(validationErrors) != 0 {
			log.Debugf("validation errors while creating acccount %+v", validationErrors)
			return web.BadRequestResponseWithDetails(c, "Niektóre pola są błędnie wypełnione", validationErrors)
		}

		accID, err := service.StartSignupAccount(account.Email, *secAccount)
		if err != nil {
			return web.LogAndReturnInternalError(c, "Could not create account ", err)
		}

		account.Id = accID
		log.Debugf("Recived ", account)
		return web.CreatedResponse(c, account)
	}

	getById := func(c echo.Context) error {

		id := c.Param("id")
		account, err := service.GetByUsername(id)

		if err != nil {

			if err == ErrAccountNotFound {
				web.NotFoundResponse(c)
			}

			return web.LogAndReturnInternalError(c, "Could not fetch accounts", err)
		}

		return c.JSON(http.StatusOK, account)
	}

	resetPassword := func(c echo.Context) error {

		email := c.Param("id")
		if err := service.StartResetPassword(email); err != nil {

			if err == ErrAccountNotFound {
				return web.NotFoundResponse(c)
			}

			return web.LogAndReturnInternalError(c, "Could not start reset password process.", err)
		}

		return c.JSON(http.StatusOK, "")
	}

	confirmResetPassword := func(c echo.Context) error {

		email := c.Param("id")
		passwordChangeDto := PasswordChangeDto{}

		if err := c.Bind(&passwordChangeDto); err != nil {
			log.Error("[confirmResetPassword]Unable to parse for confirmResetPassword", err)
			return web.BadRequestResponse(c, "Unable to parse request body")
		}

		if !passwordChangeDto.NewPassword.IsValid() {
			details := []web.ErrorDetails{{
				Field:   "password",
				Type:    web.InvalidField,
				Message: "Password is to weak",
			}}
			return web.BadRequestResponseWithDetails(c, "Password is to weak", details)
		}

		if err := service.ConfirmResetPassword(email, passwordChangeDto.Code, passwordChangeDto.NewPassword); err != nil {

			if err == ErrInvalidResetCode {
				return web.BadRequestResponse(c, "Invalid reset code")
			}

			return web.LogAndReturnInternalError(c, "Could not reset password account.", err)
		}

		return c.JSON(http.StatusOK, "")
	}

	confirmAccount := func(c echo.Context) error {

		email, _ := url.QueryUnescape(c.Param("id"))

		code := c.QueryParam("code")

		if len(code) == 0 {
			return web.BadRequestResponse(c, "Missing confirmation code")
		}

		valid, error := service.ConfirmAccount(email, code)
		if error != nil {
			return web.LogAndReturnInternalError(c, "Could not confirm account.", error)
		}

		if !valid {
			return web.BadRequestResponse(c, "Invalid confirmation code")
		}

		return c.JSON(http.StatusOK, "")

	}

	update := func(c echo.Context) error {

		acc := UpdateAccountDto{}
		if err := c.Bind(&acc); err != nil {

			log.Error("Unable to parse for update", err)
			return web.BadRequestResponse(c, "Unable to parse request body")
		}

		email, _ := url.QueryUnescape(c.Param("id"))

		if err := service.UpdateByEmail(email, acc); err != nil {

			return web.LogAndReturnInternalError(c, "unable to update account", err)
		}

		return c.JSON(http.StatusOK, "")
	}

	return Controller{
		Create:               create,
		GetByID:              getById,
		ConfirmResetPassword: confirmResetPassword,
		ConfirmAccount:       confirmAccount,
		ResetPassword:        resetPassword,
		Update:               update,
	}
}
