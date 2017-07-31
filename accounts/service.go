package accounts

import (
	"bytes"
	"github.com/op/go-logging"
	"github.com/satori/go.uuid"
	"github.com/piotrjaromin/login-template/config"
	"github.com/piotrjaromin/login-template/email"
	"github.com/piotrjaromin/login-template/dal"
)

type Service struct {
	StartSignupAccount   func(email string, secAccount SecuredAccount) (string, error)
	GetByEmail           func(email string) (PasswordlessAccount, error)
	ConfirmAccount       func(email string, code string) (bool, error)
	StartResetPassword   func(email string) error
	ConfirmResetPassword func(email string, code string, newPassword Password) error
	CreateAccount        func(email string, secAccount SecuredAccount) (string, error)
	UpdateByEmail        func(email string, accUpdate UpdateAccountDto) error
}

func CreateService(config config.Config, accountDal Dal, signupsDal dal.Dal, emailService email.EmailService, encrypt Encrypt) Service {

	var log = logging.MustGetLogger("[LoginSerivce]")
	templates := emailService.Templates()

	sendAccountRequestedMail := func(email string, code string, name string) error {

		log.Info("Signup sending signup", email)

		data := struct {
			Code  string
			Name  string
			Url   string
			Email string
		}{
			code, name, config.FrontendURL, email,
		}

		buf := new(bytes.Buffer)
		if err := templates.ExecuteTemplate(buf, "confirm_account.html", data); err != nil {
			log.Error("Cannot send confirm account email. ", err)
			return err
		}

		return emailService.SendEmail(email, buf.String(), "Account confirmation")
	}

	createAccount := func(email string, secAccount SecuredAccount) (string, error) {

		hash, salt := encrypt.Hash(secAccount.Password)
		secAccount.Password = hash
		secAccount.Salt = salt
		secAccount.Email = email
		return accountDal.CreateAccount(secAccount)
	}

	startSignup := func(email string, secAccount SecuredAccount) (string, error) {

		code := uuid.NewV4().String()

		//TODO service for signups?
		isDup, signSaveErr := signupsDal.Save(Signup{email, code})
		if signSaveErr != nil {
			log.Error("Could not save signup for ", email, ", details ", signSaveErr.Error())
		}

		if isDup {
			sign := Signup{}
			signupsDal.GetById(email, &sign)
			code = sign.Code
		}

		if err := sendAccountRequestedMail(email, code, secAccount.FirstName); err != nil {
			return "", err
		}

		secAccount.Account.Status = Pending
		return createAccount(email, secAccount)
	}

	getByEmailPasswordless := func(email string) (PasswordlessAccount, error) {
		return accountDal.GetByEmail(email)
	}

	confirmAccount := func(email string, code string) (bool, error) {

		signup := Signup{}
		if err := signupsDal.GetById(email, &signup); err != nil {
			return false, err
		}

		if signup.Code != code {
			return false, nil
		}

		if err := accountDal.UpdateByEmail(email, func(acc *SecuredAccount) error {
			acc.Status = Confirmed
			return nil
		}); err != nil {
			return false, err
		}

		signupsDal.DeleteById(email)
		return true, nil

	}

	startResetPassword := func(email string) error {

		resetCode := uuid.NewV4().String()

		updateErr := accountDal.UpdateByEmail(email, func(acc *SecuredAccount) error {
			acc.ResetPasswordCode = resetCode
			return nil
		})

		if updateErr != nil {
			log.Error("[startResetPassword] updateErr: ", updateErr.Error())
			return ErrUnableToSetResetCode
		}

		data := struct {
			Code  string
			Url   string
			Email string
		}{
			resetCode, config.FrontendURL, email,
		}

		buf := new(bytes.Buffer)
		if err := templates.ExecuteTemplate(buf, "reset_password.html", data); err != nil {
			return err
		}

		return emailService.SendEmail(email, buf.String(), "Password Reset")
	}

	confirmResetPassword := func(email string, code string, newPassword Password) error {

		handleUpdate := func(secAccount *SecuredAccount) error {
			if secAccount.ResetPasswordCode != code {
				return ErrInvalidResetCode
			}

			hash, salt := encrypt.Hash(newPassword)
			secAccount.Password = hash
			secAccount.Salt = salt
			secAccount.ResetPasswordCode = ""
			return nil
		}

		if err := accountDal.UpdateByEmail(email, handleUpdate); err != nil {
			log.Error("Error while updating account for reset password. ", err.Error())
			return err
		}

		return nil
	}

	updateByEmail := func(email string, accUpdate UpdateAccountDto) error {

		handleUpdate := func(secAccount *SecuredAccount) error {

			secAccount.Email = accUpdate.Email
			secAccount.FirstName = accUpdate.FirstName
			secAccount.LastName = accUpdate.LastName
			return nil
		}

		return accountDal.UpdateByEmail(email, handleUpdate)
	}
	return Service{
		StartSignupAccount:   startSignup,
		GetByEmail:           getByEmailPasswordless,
		ConfirmAccount:       confirmAccount,
		StartResetPassword:   startResetPassword,
		ConfirmResetPassword: confirmResetPassword,
		CreateAccount:        createAccount,
		UpdateByEmail:        updateByEmail,
	}
}
