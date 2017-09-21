package fbLogin

import (
        fb "github.com/huandu/facebook"
        "errors"
        "github.com/op/go-logging"
	"github.com/piotrjaromin/go-login-backend/accounts"
        "github.com/piotrjaromin/go-login-backend/jwtTokens"
        
)

//Errors that can be returned by this module
var (
	ErrInvalidFbToken = errors.New("Invalid fb session token was sent")
	ErrFbFetchFailed = errors.New("Could not fetch data from facebook")
	ErrCouldNotCreateAccount = errors.New("Could not create account")
	ErrCouldNotUpdateAccount = errors.New("Could not update account")
        ErrCouldNotFetchAccount = errors.New("Could not fetch account")
        ErrCouldNotGenerateToken = errors.New("Could not generate token")
)


//FbConfig with clientId and clientSecret from facebook developers site
type FbConfig struct {
        ClientID string
	ClientSecret string
}

//Service used to login with facebook
type Service struct {
	Login func(fbToken Token) (*Token, error)
}

//CreateService for fb login
func CreateService(fbConfig FbConfig, accountsDal accounts.Dal, accountsService accounts.Service, tokenService jwtTokens.TokenService) Service{

        var log = logging.MustGetLogger("[FbLoginService]")
        app := fb.New(fbConfig.ClientID, fbConfig.ClientSecret)

	updateAccount := func(fbEmail, fbId, fbFirstName, fbLastName string) error {

                return accountsDal.UpdateByEmail(fbEmail, func(acc *accounts.SecuredAccount) error {

                        log.Infof("Acc is %+v\n", *acc)
                        acc.AuthProviders.FB = fbId
                        if len(acc.FirstName) == 0 {
                                log.Info("update first name ", fbFirstName)
                                acc.FirstName = fbFirstName
                        }

                        if len(acc.LastName) == 0 {
                                log.Info("update last name ", fbLastName)
                                acc.LastName = fbLastName
                        }
                        return nil
                })
        }

	login := func(fbToken Token) (*Token, error) {
		
                session := app.Session(fbToken.Token)

                if err := session.Validate(); err != nil {
                        log.Warning("Invalid session token was sent: " + err.Error())
                        return nil, ErrInvalidFbToken
                }

                res, fbErr := session.Get("/me", fb.Params{
                        "fields":       "first_name,last_name,email,id",
                        "access_token": fbToken.Token,
                })

                if fbErr != nil {
                        return nil, ErrFbFetchFailed
                }

                fbEmail := res["email"].(string)
                fbID := res["id"].(string)
                fbFirstName := res["first_name"].(string)
                fbLastName := res["last_name"].(string)

                acc, err := accountsDal.GetByEmail(fbEmail)

                accountID := acc.Id
                if err == accounts.ErrAccountNotFound {

                        secAcc := accounts.SecuredAccount{
                                Account: accounts.Account{
                                        PasswordlessAccount: accounts.PasswordlessAccount{
                                                Id: fbEmail,
                                                FirstName: fbFirstName,
                                                LastName: fbLastName,
                                                AuthProviders: accounts.AuthProviders{FB:fbID},
                                                Username: fbToken.Username,
                                        },
                                },
                        }
                        if accountID, err = accountsService.CreateAccount(fbEmail, secAcc); err != nil {
                                return nil, ErrCouldNotCreateAccount
                        }

                } else if err != nil {
                        return nil, ErrCouldNotFetchAccount
                }

                if acc.AuthProviders.FB != fbID {
                        if err := updateAccount(fbEmail, fbID, fbFirstName, fbLastName); err != nil {
                                 return nil, ErrCouldNotUpdateAccount
                        }                       
                }

                tokenStr, err := tokenService.GenerateToken(fbEmail, accountID)
                if err != nil {
                        return nil, ErrCouldNotGenerateToken
                }

                return &Token{
                        Token: tokenStr,
                }, nil
	}

	return Service{
		Login: login,
	}
}