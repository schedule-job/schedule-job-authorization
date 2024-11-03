package core

type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Provider struct {
	Name     string `json:"name"`
	LoginUrl string `json:"loginUrl"`
}

type OAuth struct {
	providers map[string]OAuthInterface
}
