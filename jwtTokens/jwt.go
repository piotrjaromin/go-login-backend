package jwtTokens

import (
        "time"
        "github.com/dgrijalva/jwt-go"
        "github.com/op/go-logging"
        "fmt"
)

type TokenService struct {
        Validate      func(token string, claimName string, claimValue string) bool
        GenerateToken func(username string, userId string) (string, error)
}

//Create service which generates and validates jwt tokens
func Create(signingKey string) TokenService {
        var log = logging.MustGetLogger("[jwtTokens]")

        validate := func(tokenString string, claimName string, claimValue string) bool {

                token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                        // Don't forget to validate the alg is what you expect:
                        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                                return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
                        }
                        return []byte(signingKey), nil
                })

                if token == nil {
                        log.Debugf("Token is nil %+v", tokenString)
                        return false
                }

                if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
                        log.Debugf("Validating claim %s against value %s, token contiains %s", claimName, claimValue, claims[claimName])
                        return claims[claimName] == claimValue
                }

                log.Errorf("Error while validating token %+v", err)
                return false
        }

        generateToken := func(username string, id string) (string, error) {
                token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
                        "username": username,
                        "userId": id,
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
                Validate:validate,
                GenerateToken:generateToken,
        }
}


