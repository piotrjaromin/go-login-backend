package main

import (
        "github.com/labstack/echo"
        "github.com/op/go-logging"
        "github.com/labstack/echo/engine/standard"
        "github.com/labstack/echo/middleware"

        "github.com/piotrjaromin/go-login-backend/jwtTokens"
        "github.com/piotrjaromin/go-login-backend/config"
        "github.com/piotrjaromin/go-login-backend/security"
        "github.com/piotrjaromin/go-login-backend/accounts"
        "github.com/piotrjaromin/go-login-backend/login"
        "github.com/piotrjaromin/go-login-backend/fbLogin"
        "github.com/piotrjaromin/go-login-backend/email"
        "github.com/piotrjaromin/go-login-backend/dal"
)

var log = logging.MustGetLogger("[Main]")

func main() {

        conf := config.GetConfig("./config/" + config.GetEnvOrDefault("CONF_FILE", "config.json"))

        log.Info("Preparing app")
        e := echo.New()

        e.Use(middleware.Logger())
        e.Use(middleware.Recover())

        headers := func(next echo.HandlerFunc) echo.HandlerFunc {
                return func(c echo.Context) error {
                        c.Response().Header().Add("Content-type", "application/json")
                        c.Response().Header().Add("Allow", "GET,POST,HEAD,OPTIONS,PUT")
                        c.Response().Header().Add("Access-Control-Allow-Methods", "GET,POST,HEAD,OPTIONS,PUT")
                        c.Response().Header().Add("Access-Control-Allow-Origin", "*");
                        c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With");
                        c.Response().Header().Add("Access-Control-Max-Age", "3600");
                        return next(c)
                }
        }

        e.Use(headers)

        tokenService := jwtTokens.Create(conf.Token.SiginKey)
        security := security.CreateSecurity(tokenService)

        //Accounts endpoints
        accDal := accounts.CreateDal(getCollection("accounts", conf))
        singupDal := getCollection("signups", conf)
        emailService, emailErr := email.Create(conf.Email.AwsRegion, conf.Email.ReplyAddr)
        if emailErr != nil {
                panic("Could not creat email service. Details: " + emailErr.Error())
        }

        encrypt := accounts.CreateEncrypt()
        accService := accounts.CreateService(conf, accDal, singupDal, emailService, encrypt)
        accController := accounts.Create(accService)
        accounts.InitRoutes(e, accController, security)        

        //Login endpoints
        loginService := login.CreateService(accDal, encrypt, tokenService)
        loginController := login.Create(loginService)
        login.InitRoutes(e, loginController)

        //Fb login endpoints
        fbLoginService := fbLogin.CreateService(fbLogin.FbConfig{
                ClientID: conf.Fb.ClientID,
                ClientSecret: conf.Fb.ClientSecret,
        }, accDal, accService, tokenService)
        
        fbLoginController := fbLogin.Create(fbLoginService)
        fbLogin.InitRoutes(e, fbLoginController)

        createAccount(accService)

        log.Info("Starting to listen")
        error := e.Run(standard.New(":8080"))

        if error != nil {
                log.Errorf("Error while starting server %+v", error)
        }
}

func getCollection(collection string, conf config.Config) dal.Dal {
        config := dal.DalConfig{
                Server: conf.Mongo.Server,
                Database: conf.Mongo.Database,
                Collection: collection,
        }

        return dal.Create(config)
}


func createAccount(accService accounts.Service) {

        const email = "test@test.com"

        _, err := accService.GetByEmail(email)

        //Create test account if it does not exist
        if err != nil && err == accounts.ErrAccountNotFound {
                log.Info("Creating test account")

                accService.CreateAccount(email, accounts.SecuredAccount{
                        Account: accounts.Account{
                                Password: "test",
                                PasswordlessAccount: accounts.PasswordlessAccount{
                                        Email: email,
                                        FirstName: "fn",
                                        LastName: "ln",
                                        Status: accounts.Confirmed,
                                },
                        },
                })
        }

}