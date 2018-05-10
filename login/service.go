package login

import (
	"errors"
	"github.com/op/go-logging"
	"github.com/piotrjaromin/go-login-backend/accounts"
	"github.com/piotrjaromin/go-login-backend/jwtTokens"
)

//Errors that can be returned by this module
var (
	ErrMissingPasswordOrUsername = errors.New("Both password and username are required")
	ErrCouldNotFetchAccount      = errors.New("Error while fetching account for authentication")
	ErrNotFoundAccount           = errors.New("Account does not exist")
	ErrNotConfirmedAccount       = errors.New("Account is not confirmed")
	ErrBadCredentials            = errors.New("Invalid credentials")
	ErrCouldNotGenerateToken     = errors.New("Could not generate token")
)

//Service with login function
type Service struct {
	Login func(username string, pass accounts.Password) (*Token, error)
}

//CreateService creates service responsible for issuing tokens
func CreateService(accountsDal accounts.Dal, encrypt accounts.Encrypt, tokenService jwtTokens.TokenService) Service {

	var log = logging.MustGetLogger("[LoginService]")

	login := func(username string, pass accounts.Password) (*Token, error) {

		log.Debug("got from query params ", username, pass)
		if len(username) == 0 || len(pass) == 0 {
			return nil, ErrMissingPasswordOrUsername
		}

		secAccount, getAccErr := accountsDal.GetWithPasswordByEmail(username)
		if getAccErr != nil {
			if getAccErr == accounts.ErrAccountNotFound {
				return nil, ErrBadCredentials
			}
			log.Error("Cold not fetch account. Detials: ", getAccErr.Error())
			return nil, ErrCouldNotFetchAccount
		}

		log.Debugf("found user %s", secAccount.Id)
		if secAccount.Status != accounts.Confirmed {
			return nil, ErrNotConfirmedAccount
		}

		if !encrypt.Validate(pass, secAccount.Password, secAccount.Salt) {
			return nil, ErrBadCredentials
		}

		tokenStr, err := tokenService.GenerateToken(secAccount.Username, secAccount.Id)
		if err != nil {
			return nil, ErrCouldNotGenerateToken
		}

		token := Token{
			Token: tokenStr,
		}

		log.Infof("Generated token is %s", token.Token)
		return &token, nil
	}

	return Service{
		login,
	}
}
