package accounts

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/piotrjaromin/go-login-backend/test"
	"github.com/piotrjaromin/go-login-backend/web"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestController(t *testing.T) {

	validAccount := Account{
		PasswordlessAccount: PasswordlessAccount{
			FirstName: "Jhone",
			LastName:  "Doe",
			Email:     "test@test.com",
			Username:  "testUser",
		},
		Password: "123456aA",
	}

	invalidAccount := Account{
		PasswordlessAccount: PasswordlessAccount{
			FirstName: "Jhone",
			LastName:  "Doe",
		},
		Password: "123A",
	}

	accountsService := func() Service {
		return Service{
			StartSignupAccount: func(email string, secAccount SecuredAccount) (string, error) {

				return "testID", nil
			},
			GetByUsername: func(id string) (PasswordlessAccount, error) {

				if id == validAccount.Email {
					return validAccount.PasswordlessAccount, nil
				}

				return PasswordlessAccount{}, ErrAccountNotFound
			},
			StartResetPassword: func(email string) error {
				if email == validAccount.Email {
					return nil
				}

				return ErrAccountNotFound
			},
		}
	}

	createContextAndRecorder := func(controller Controller, req *http.Request) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()

		e := echo.New()
		InitRoutes(e, controller, test.CreateSecurity("valid"))

		res := echo.NewResponse(rec, e)
		e.ServeHTTP(res, req)
		return rec
	}

	Convey("for post on accounts", t, func() {

		Convey("should start singup process", func() {

			accountJson, _ := json.Marshal(validAccount)
			req, _ := http.NewRequest(echo.POST, "/accounts", strings.NewReader(string(accountJson)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			resp := createContextAndRecorder(Create(accountsService()), req)

			So(resp.Code, ShouldEqual, http.StatusCreated)

			createdResp := Account{}

			json.Unmarshal(resp.Body.Bytes(), &createdResp)
			So(createdResp.Email, ShouldEqual, validAccount.Email)
		})

		Convey("should return bad request when invalid account data is sent", func() {

			accountJson, _ := json.Marshal(invalidAccount)
			req, _ := http.NewRequest(echo.POST, "/accounts", strings.NewReader(string(accountJson)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			resp := createContextAndRecorder(Create(accountsService()), req)

			So(resp.Code, ShouldEqual, http.StatusBadRequest)

			errorDto := web.Error{}

			json.Unmarshal(resp.Body.Bytes(), &errorDto)

			So(errorDto.ErrorDetails, ShouldHaveLength, 3)
		})
	})

	Convey("for get on single account should", t, func() {

		Convey("return existing account", func() {

			req, _ := http.NewRequest(echo.GET, "/accounts/"+validAccount.Email, strings.NewReader(""))
			req.Header.Set("Authorization", "Bearer valid")

			resp := createContextAndRecorder(Create(accountsService()), req)
			So(resp.Code, ShouldEqual, http.StatusOK)

			account := SecuredAccount{}
			json.Unmarshal(resp.Body.Bytes(), &account)

			So(account.Email, ShouldEqual, validAccount.Email)
			So(account.FirstName, ShouldEqual, validAccount.FirstName)
			So(account.LastName, ShouldEqual, validAccount.LastName)
			So(account.Password, ShouldEqual, Password(""))
			So(account.Salt, ShouldBeBlank)

		})

		Convey("return not found for non existing account", func() {

			req, _ := http.NewRequest(echo.GET, "/accounts/ranom@mail.com", strings.NewReader(""))
			req.Header.Set("Authorization", "Bearer valid")

			resp := createContextAndRecorder(Create(accountsService()), req)

			So(resp.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("return access forbidden for invalid token", func() {

			req, _ := http.NewRequest(echo.GET, "/accounts/ranom@mail.com", strings.NewReader(""))
			req.Header.Set("Authorization", "Bearer invalidToken")

			resp := createContextAndRecorder(Create(accountsService()), req)

			So(resp.Code, ShouldEqual, http.StatusUnauthorized)
		})
	})

	Convey("for post on accounts/:id/reset should", t, func() {

		Convey("start reset password flow", func() {

			req, _ := http.NewRequest(echo.POST, "/accounts/"+validAccount.Email+"/reset", strings.NewReader(""))

			resp := createContextAndRecorder(Create(accountsService()), req)
			So(resp.Code, ShouldEqual, http.StatusOK)

		})

		Convey("return not found for not existing account", func() {

			req, _ := http.NewRequest(echo.POST, "/accounts/random@mail.com/reset", strings.NewReader(""))

			resp := createContextAndRecorder(Create(accountsService()), req)
			So(resp.Code, ShouldEqual, http.StatusNotFound)

		})
	})

	Convey("for put on accounts/:id/reset should", t, func() {

		validCode := "validCode12"
		accService := Service{
			ConfirmResetPassword: func(email string, code string, newPassword Password) error {
				if email != validAccount.Email || code != validCode {
					return ErrInvalidResetCode
				}
				return nil
			},
		}

		Convey("confirm reset password", func() {

			passDto := PasswordChangeDto{
				Code:        validCode,
				NewPassword: "123456aA",
			}

			passDtoJson, _ := json.Marshal(passDto)

			req, _ := http.NewRequest(echo.PUT, "/accounts/"+validAccount.Email+"/reset", strings.NewReader(string(passDtoJson)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			resp := createContextAndRecorder(Create(accService), req)
			So(resp.Code, ShouldEqual, http.StatusOK)

		})

		Convey("return bad request for too weak password", func() {

			passDto := PasswordChangeDto{
				Code:        validCode,
				NewPassword: "6aA",
			}

			passDtoJson, _ := json.Marshal(passDto)

			req, _ := http.NewRequest(echo.PUT, "/accounts/"+validAccount.Email+"/reset", strings.NewReader(string(passDtoJson)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			resp := createContextAndRecorder(Create(accService), req)
			So(resp.Code, ShouldEqual, http.StatusBadRequest)

		})

		Convey("return bad request for invalid reset code", func() {

			passDto := PasswordChangeDto{
				Code:        "randomCode",
				NewPassword: "12345678aA",
			}

			passDtoJson, _ := json.Marshal(passDto)

			req, _ := http.NewRequest(echo.PUT, "/accounts/"+validAccount.Email+"/reset", strings.NewReader(string(passDtoJson)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			resp := createContextAndRecorder(Create(accService), req)

			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})

	})

	Convey("for put on accounts/:id/confirm should", t, func() {

		validCode := "validCode12"
		accService := Service{
			ConfirmAccount: func(email string, code string) (bool, error) {
				if email != validAccount.Email {
					return false, ErrAccountNotFound
				}

				if code != validCode {
					return false, nil
				}

				return true, nil
			},
		}

		Convey("confirm account", func() {

			req, _ := http.NewRequest(echo.GET, "/accounts/"+validAccount.Email+"/confirm?code="+validCode, strings.NewReader(string("")))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			resp := createContextAndRecorder(Create(accService), req)
			So(resp.Code, ShouldEqual, http.StatusOK)

		})

		Convey("return bad request for invalid code", func() {

			req, _ := http.NewRequest(echo.GET, "/accounts/"+validAccount.Email+"/confirm?code=invalidCode", strings.NewReader(string("")))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			resp := createContextAndRecorder(Create(accService), req)
			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("return error when error occured", func() {
			req, _ := http.NewRequest(echo.GET, "/accounts/notexistin@mail.com/confirm?code=invalidCode", strings.NewReader(string("")))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			resp := createContextAndRecorder(Create(accService), req)
			So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		})

	})

}
