package core

type OAuthInterface interface {
	GetUser(code string) (*User, error)
	GetLoginUrl() string
}
