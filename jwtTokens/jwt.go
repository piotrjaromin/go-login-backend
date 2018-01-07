package jwtTokens

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/op/go-logging"
)

type TokenService struct {
	Validate      func(token string, claimName string, claimValue string) bool
	GenerateToken func(username string, userId string) (string, error)
	GetClaims     func(tokenString string) map[string]interface{}
}

//Create service which generates and validates jwt tokens
func Create(signingKey string) TokenService {
	var log = logging.MustGetLogger("[jwtTokens]")

	getClaims := func(tokenString string) map[string]interface{} {

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(signingKey), nil
		})

		if token == nil || err != nil {
			log.Debugf("Token is nil %+v or error occured %+v", tokenString, err)
			return map[string]interface{}{}
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims
		}

		return map[string]interface{}{}
	}

	validate := func(tokenString string, claimName string, claimValue string) bool {

		claims := getClaims(tokenString)
		log.Debugf("Validating claim %s against value %s, token contiains %s", claimName, claimValue, claims[claimName])
		return claims[claimName] == claimValue
	}

	generateToken := func(username string, id string) (string, error) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username":  username,
			"userId":    id,
			"expiresAt": time.Now().Add(time.Hour * 24 * 7).Unix(),
		})

		tokenString, err := token.SignedString([]byte(signingKey))

		if err != nil {
			log.Errorf("Could not generate token string %+v", err)
			return "", err
		}

		return tokenString, nil
	}

	return TokenService{
		Validate:      validate,
		GenerateToken: generateToken,
		GetClaims:     getClaims,
	}
}
