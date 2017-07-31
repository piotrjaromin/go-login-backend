package web

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("[web]")

func CreatedResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, data)
}

func UnauthorizedResponse(c echo.Context, msg string) error {

	resp := Error{
		Message: msg,
		Status:  http.StatusUnauthorized,
	}

	return c.JSON(http.StatusUnauthorized, resp)
}

func NotFoundResponse(c echo.Context) error {

	resp := Error{
		Message: "Resource does not exist",
		Status:  http.StatusNotFound,
	}

	return c.JSON(http.StatusNotFound, resp)
}

func BadRequestResponse(c echo.Context, msg string) error {

	resp := Error{
		Message: msg,
		Status:  http.StatusBadRequest,
	}

	return c.JSON(http.StatusBadRequest, resp)
}

func ConflictResponse(c echo.Context, msg string) error {

	resp := Error{
		Message: msg,
		Status:  http.StatusConflict,
	}

	return c.JSON(http.StatusConflict, resp)
}


func BadRequestResponseWithDetails(c echo.Context, msg string, details []ErrorDetails) error {

	resp := Error{
		Message:      msg,
		Status:       http.StatusBadRequest,
		ErrorDetails: details,
	}

	return c.JSON(http.StatusBadRequest, resp)
}

func LogAndReturnInternalError(c echo.Context, msg string, err error) error {

	logger.Errorf(msg, ". Details: ", err)
	resp := Error{
		Message: msg,
		Status:  http.StatusInternalServerError,
	}

	return c.JSON(http.StatusInternalServerError, resp)
}
