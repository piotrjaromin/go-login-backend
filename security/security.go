package security

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"regexp"

	"github.com/labstack/echo"
	"github.com/op/go-logging"
	"github.com/piotrjaromin/go-login-backend/jwtTokens"
	"github.com/piotrjaromin/go-login-backend/web"
)

const CLAIM_IN_QUOTES_PATTERN = "\".*\""

type Security struct {
	tokenService jwtTokens.TokenService
}

func CreateSecurity(tokenService jwtTokens.TokenService) Security {
	return Security{tokenService: tokenService}
}

func (sec Security) SecuredById(tokenClaimName string, requestClaimName string, claimInBody bool) func(next echo.HandlerFunc) echo.HandlerFunc {

	var log = logging.MustGetLogger("[Security]")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if c.Request().Method() == "OPTIONS" {
				return nil
			}

			token, found := getToken(c)
			if !found {
				return web.UnauthorizedResponse(c, "Invalid authorization header")
			}

			var idValue string
			if claimInBody {
				//read idCliamName from body and check if it matches the one stored in token
				var jsonMap map[string]*json.RawMessage
				var data, reqErr = ioutil.ReadAll(c.Request().Body())
				if reqErr != nil {
					return web.LogAndReturnInternalError(c, "[sec]error while reading request body", reqErr)
				}

				c.Request().SetBody(ioutil.NopCloser(bytes.NewBuffer(data)))

				err := json.Unmarshal(data, &jsonMap)
				if err != nil {
					return web.LogAndReturnInternalError(c, "[sec]error while unmarshalling request body", err)
				}

				claimValue := jsonMap[requestClaimName]
				if claimValue == nil {
					return web.BadRequestResponse(c, "body does not contain required claim "+requestClaimName)
				}

				idValue = string(*claimValue)

				//string in body contains " at begining and end, remove them
				claimInQuotes, _ := regexp.MatchString(CLAIM_IN_QUOTES_PATTERN, idValue)
				if claimInQuotes {
					idValue = idValue[1 : len(idValue)-1]
				}
			} else {
				idValue = c.Param("id")
			}

			valid := sec.tokenService.Validate(token, tokenClaimName, idValue)

			if valid {
				log.Info("saving " + tokenClaimName + " with value " + idValue)
				c.Set(tokenClaimName, idValue)
				return next(c)
			}

			log.Info("invalid claim")
			return web.UnauthorizedResponse(c, "You do not have required scopes to perform this method")
		}
	}
}

func (sec Security) FillClaims() func(next echo.HandlerFunc) echo.HandlerFunc {
	var log = logging.MustGetLogger("[Security]")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			token, found := getToken(c)
			if found {
				claims := sec.tokenService.GetClaims(token)
				for key, value := range claims {
					c.Set(key, value)
				}
			}

			return next(c)
		}
	}
}

func getToken(c echo.Context) (string, bool) {
	authHeader := c.Request().Header().Get("Authorization")

	if len(authHeader) < len("Bearer ") {
		return "", false
	}

	return authHeader[len("Bearer "):], true
}
