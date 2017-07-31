package accounts

import (
	"regexp"
	"time"
	"unicode"
	e "github.com/piotrjaromin/login-template/web"
	"errors"
)

type AccountStatus string

//Errors returned by this module
var (
	ErrInvalidResetCode = errors.New("Invalid reset password code")
	ErrUnableToSetResetCode = errors.New("Invalid reset password code")
	ErrAccountNotFound = errors.New("Account does not exist")
)

const (
	Pending   AccountStatus = "PENDING"
	Confirmed AccountStatus = "CONFIRMED"
)

type Password string

func (p Password) IsValid() bool {

	fiveOrMore, number, upper := verifyPassword(p)
	if !fiveOrMore || !number || !upper {
		return false
	}
	return true
}

type UpdateAccountDto struct {
	Email          string `json:"email"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
}

type PasswordlessAccount struct {
	Id             string    `bson:"_id"`
	Email          string    `json:"email" bson:"email"`
	FirstName      string    `json:"firstName" bson:"firstName"`
	LastName       string    `json:"lastName" bson:"lastName"`
	CreatedAt      time.Time `json:"createdAt,omitempty" bson:"createdAt"`
	Status         AccountStatus `bson:"status"`
	AuthProviders  AuthProviders `bson:"authProviders"`
}

type PasswordChangeDto struct {
	Code        string   `json:"code"`
	NewPassword Password `json:"newPassword"`
}

type ConfirmAccountDto struct {
	Code string `json:"code"`
}

type Account struct {
	PasswordlessAccount `bson:",inline"`
	Password            Password `json:"password,omitempty" bson:"password"`
}

type Signup struct {
	Email string `json:"email" bson:"_id"`
	Code  string `json:"code" bson:"code"`
}

type AuthProviders struct {
	FB string `bson:"FB"`
}

type SecuredAccount struct {
	Account           `bson:",inline"`
	Salt              string `json:"salt" bson:"salt"`
	ResetPasswordCode string `bson:"resetPasswordCode"`
}

const (
	emailPattern = "[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,6}"
)

func (acc Account) validate() []e.ErrorDetails {

	var errors []e.ErrorDetails

	if len(acc.FirstName) == 0 {
		errors = e.AppendErrorDetails(errors, "firstName", "First name is required", e.MissingField)
	}

	if len(acc.LastName) == 0 {
		errors = e.AppendErrorDetails(errors, "lastName", "Last name is required", e.MissingField)
	}

	if !acc.Password.IsValid() {
		errors = e.AppendErrorDetails(errors, "password", "Password to weak(5 characters, one capital letter, one digit)", e.InvalidField)
	}

	emailOk, _ := regexp.MatchString(emailPattern, acc.Email)
	if !emailOk {
		errors = e.AppendErrorDetails(errors, "email", "Email in invalid format", e.InvalidField)
	}

	return errors
}

func verifyPassword(s Password) (fiveOrMore, number, upper bool) {

	pass := string(s)
	for _, s := range pass {
		switch {
		case unicode.IsNumber(s):
			number = true
		case unicode.IsUpper(s):
			upper = true
		}
	}
	fiveOrMore = len(pass) >= 5
	return
}
