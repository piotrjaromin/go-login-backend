package test

import (
	"github.com/piotrjaromin/go-login-backend/config"
	"github.com/piotrjaromin/go-login-backend/jwtTokens"
	"github.com/piotrjaromin/go-login-backend/dal"
	"github.com/piotrjaromin/go-login-backend/security"
)

func GetTestRepo(collection string) dal.Dal {

	conf := config.GetConfig("./config/" + config.GetEnvOrDefault("CONF_FILE", "test.json"))

	return dal.Create(dal.DalConfig{
		Server:     conf.Mongo.Server,
		Database:   conf.Mongo.Database,
		Collection: collection,
	})
}

func CreateSecurity(validToken string) security.Security {

	tokenService := jwtTokens.TokenService{
		Validate: func(token string, claimName string, claimValue string) bool {
			return token == validToken
		},
		GenerateToken: func(username string, userId string) (string, error) {

			return validToken, nil
		},
	}

	return security.CreateSecurity(tokenService)
}
