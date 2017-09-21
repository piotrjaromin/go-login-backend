package accounts

import (
        . "github.com/smartystreets/goconvey/convey"
        "testing"
        "github.com/piotrjaromin/go-login-backend/config"
        "github.com/piotrjaromin/go-login-backend/dal"
        "github.com/smartystreets/assertions/should"
        "html/template"
        "github.com/satori/go.uuid"
)

type TestMail struct {
        expected string
}

func (t TestMail) SendEmail(mail string, content string, subject string) error {
        So(mail, should.Equal, t.expected)
        return nil
}

func (t TestMail) Templates() *template.Template {
        tmp, _ := template.New("confirm_account.html").Parse("test confirm template")
        return template.Must(tmp.New("reset_password.html").Parse("test reset template"))
}

func TestService(t *testing.T) {

        testEmail := "test@test.com"

        validSecuredAcc := SecuredAccount{
                Account: Account{
                        PasswordlessAccount: PasswordlessAccount{
                                FirstName : "firstName from test",
                                LastName : "lastName from test",
                                Email: testEmail,
                        },
                        Password : "testPass",
                },
        }

        Convey("CreateAccount should", t, func() {

                testSalt := "testSalt"
                hashedPass := Password("hashedPassword")

                accountsRepo := Dal{
                        CreateAccount: func(secAccount SecuredAccount) (string, error) {

                                So(secAccount.Email, should.Equal, validSecuredAcc.Email)
                                So(secAccount.FirstName, should.Equal, validSecuredAcc.FirstName)
                                So(secAccount.LastName, should.Equal, validSecuredAcc.LastName)
                                So(secAccount.Status, should.Equal, Pending)
                                So(secAccount.Password, should.Equal, hashedPass)
                                So(secAccount.Salt, should.Equal, testSalt)
                                return uuid.NewV4().String(), nil
                        },
                }

                signupsRepo := dal.Dal{
                        Save: func(data interface{}) (bool, error) {

                                signup, ok := data.(Signup)
                                So(ok, ShouldBeTrue)
                                So(signup.Email, should.Equal, testEmail)
                                So(signup.Code, should.NotBeBlank)
                                return false, nil
                        },
                }

                emailService := TestMail{testEmail}
                encrypt := Encrypt{
                        Hash: func(pass Password) (Password, string) {
                                return hashedPass, testSalt
                        },
                }
                service := CreateService(config.Config{}, accountsRepo, signupsRepo, emailService, encrypt)

                Convey("Create valid account", func() {

                        _, err := service.StartSignupAccount(testEmail, validSecuredAcc)

                        So(err, should.BeNil)
                })
        })

        Convey("GetById should", t, func() {

                signupsRepo := dal.Dal{}
                emailService := TestMail{}
                encrypt := Encrypt{}

                Convey("return account", func() {

                        accountsDal := Dal{
                                GetByEmail: func(email string) (PasswordlessAccount, error) {
                                        So(email, should.Equal, testEmail)

                                        acc := PasswordlessAccount{
                                                Id: uuid.NewV4().String(),
                                                Email: email,
                                                FirstName: validSecuredAcc.FirstName,
                                                LastName: validSecuredAcc.LastName,
                                        }
                                        return acc, nil
                                },
                        }

                        service := CreateService(config.Config{}, accountsDal, signupsRepo, emailService, encrypt)

                        acc, error := service.GetByEmail(testEmail)

                        So(error, should.BeNil)

                        So(acc.FirstName, should.Equal, validSecuredAcc.FirstName)
                        So(acc.LastName, should.Equal, validSecuredAcc.LastName)

                })

                Convey("return found false when account does not exist", func() {

                        accountsDal := Dal{
                                GetByEmail: func(email string) (PasswordlessAccount, error) {
                                        So(email, should.Equal, testEmail)
                                        return PasswordlessAccount{}, ErrAccountNotFound
                                },
                        }

                        service := CreateService(config.Config{}, accountsDal, signupsRepo, emailService, encrypt)

                        _, error := service.GetByEmail(testEmail)

                        So(error, should.Equal, ErrAccountNotFound)
                })
        })

        Convey("ConfirmAccount should", t, func() {

                confirmCode := "testCode"

                accountDal := Dal{
                        UpdateByEmail: func(email string, handleUpdateFunc func(*SecuredAccount) error) (error) {
                                So(email, should.Equal, testEmail)

                                secAcc := SecuredAccount{}
                                handleUpdateFunc(&secAcc)
                                So(secAcc.Status, should.Equal, Confirmed)
                                return nil
                        },
                }

                signupsRepo := dal.Dal{
                        GetById: func(id string, data interface{}) (error) {

                                signup, ok := data.(*Signup)
                                So(ok, ShouldBeTrue)
                                signup.Code = confirmCode
                                signup.Email = id
                                return nil
                        },
                        DeleteById:    func(id string) error {
                                So(id, should.Equal, testEmail)
                                return nil;
                        },
                }

                emailService := TestMail{}
                encrypt := Encrypt{}

                Convey("change status of account to confirmed if code is valid", func() {

                        service := CreateService(config.Config{}, accountDal, signupsRepo, emailService, encrypt)

                        ok, err := service.ConfirmAccount(testEmail, confirmCode)

                        So(err, should.BeNil)
                        So(ok, should.BeTrue)
                })

                Convey("return invalid code error for wrong code", func() {

                        service := CreateService(config.Config{}, accountDal, signupsRepo, emailService, encrypt)

                        ok, err := service.ConfirmAccount(testEmail, "InvalidCode")

                        So(err, should.BeNil)
                        So(ok, should.BeFalse)
                })
        })

        Convey("ConfirmResetPassword should", t, func() {

                confirmCode := "testCode"
                newPass := Password("newPassword123!")
                hashedPass := Password("testHashed")
                testSalt := "testSalt"

                signupsRepo := dal.Dal{}
                emailService := TestMail{}
                encrypt := Encrypt{
                        Hash: func(pass Password) (Password, string) {
                                return hashedPass, testSalt
                        },
                }

                Convey("change password of account if code is valid", func() {

                        accountDal := Dal{
                                UpdateByEmail: func(email string, handleUpdateFunc func(*SecuredAccount) error) (error) {
                                        So(email, should.Equal, testEmail)

                                        secAcc := SecuredAccount{}
                                        secAcc.ResetPasswordCode = confirmCode
                                        err := handleUpdateFunc(&secAcc)
                                        So(err, ShouldBeNil)
                                        So(secAcc.Password, should.Equal, hashedPass)
                                        So(secAcc.Salt, should.Equal, testSalt)
                                        So(secAcc.ResetPasswordCode, should.BeBlank)
                                        return nil
                                },
                        }

                        service := CreateService(config.Config{}, accountDal, signupsRepo, emailService, encrypt)

                        err := service.ConfirmResetPassword(testEmail, confirmCode, newPass)

                        So(err, should.BeNil)
                })

                Convey("return invalid code error for wrong code", func() {

                        accountDal := Dal{
                                UpdateByEmail: func(email string, handleUpdateFunc func(*SecuredAccount) error) (error) {
                                        So(email, should.Equal, testEmail)
                                        secAcc := SecuredAccount{}
                                        err := handleUpdateFunc(&secAcc)
                                        return err
                                },
                        }

                        service := CreateService(config.Config{}, accountDal, signupsRepo, emailService, encrypt)

                        err := service.ConfirmResetPassword(testEmail, "InvalidCode", newPass)
                        So(err, should.Equal, ErrInvalidResetCode)
                })

        })

        Convey("StartResetPassword should", t, func() {

                accountsDal := Dal{
                        UpdateByEmail: func(email string, handleUpdateFunc func(*SecuredAccount) error) (error) {
                                So(email, should.Equal, testEmail)

                                secAcc := SecuredAccount{}
                                err := handleUpdateFunc(&secAcc)

                                So(secAcc.ResetPasswordCode, ShouldNotBeBlank)

                                return err
                        },
                }

                signupsRepo := dal.Dal{};
                emailService := TestMail{testEmail}
                encrypt := Encrypt{}

                Convey("should add reset code", func() {

                        service := CreateService(config.Config{}, accountsDal, signupsRepo, emailService, encrypt)

                        err := service.StartResetPassword(testEmail)

                        So(err, should.BeNil)
                })

                Convey("should return not found error for not existing account", func() {
                        accountsDal := Dal{
                                UpdateByEmail: func(email string, handleUpdateFunc func(*SecuredAccount) error) (error) {
                                        return ErrAccountNotFound
                                },
                        }

                        service := CreateService(config.Config{}, accountsDal, signupsRepo, emailService, encrypt)

                        err := service.StartResetPassword("not@existing.com")

                        So(err, should.Equal, ErrUnableToSetResetCode)
                })
        })
}