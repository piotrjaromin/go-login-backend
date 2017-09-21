package security

import (
	"github.com/piotrjaromin/go-login-backend/jwtTokens"
	"github.com/piotrjaromin/go-login-backend/web"
	"github.com/labstack/echo"
	"github.com/op/go-logging"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"regexp"
)

const CLAIM_IN_QUOTES_PATTERN = "\".*\""

type Security struct {
	tokenService jwtTokens.TokenService
}

func CreateSecurity(tokenService jwtTokens.TokenService) Security {
	return Security{tokenService: tokenService}
}


func (sec Security) SecuredById(tokenClaimName string, requestClaimName string, claimInBody bool) (func(next echo.HandlerFunc) echo.HandlerFunc) {

	var log = logging.MustGetLogger("[Security]")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if c.Request().Method() == "OPTIONS" {
				return nil;
			}

			authHeader := c.Request().Header().Get("Authorization")

			if len(authHeader) < len("Bearer ") {
				return web.UnauthorizedResponse(c, "Invalid authorization header")
			}

			token := authHeader[len("Bearer "):]

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

				idValue = string(*jsonMap[requestClaimName])

				//string in body contains " at begining and end, remove them
				claimInQuotes, _ := regexp.MatchString(CLAIM_IN_QUOTES_PATTERN, idValue)
				if claimInQuotes {
					idValue = idValue[1:len(idValue) - 1]
				}
			} else {
				idValue = c.Param("id")
			}

			valid := sec.tokenService.Validate(token, tokenClaimName, idValue)

			if valid {
				log.Info("saving " + tokenClaimName + " with value " + idValue)
				c.Set(tokenClaimName, idValue)
				return next(c)
			} else {
				log.Info("invalid claim")
				return web.UnauthorizedResponse(c, "You do not have required scopes to perform this method")
			}
		}
	}
}