package fbLogin

import (
        e "github.com/piotrjaromin/go-login-backend/web"
)

//Token model for fb token
type Token struct {
        Token string `json:"token"`
}


func (token Token) Validate() (errors []e.ErrorDetails) {

        if len(token.Token) == 0 {
		errors = e.AppendErrorDetails(errors, "token", "token is required", e.MissingField)
        }

        return
}