package fbLogin

//Token model for fb token
type Token struct {
        Token string `json:"token"`
}

//TokenWithEmail model
type TokenWithEmail struct {
        Token
        Email string `json:"email"`
}