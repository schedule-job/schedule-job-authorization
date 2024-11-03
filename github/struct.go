package github

import "github.com/schedule-job/schedule-job-authorization/core"

type Github struct {
	ClientId             string
	ClientSecret         string
	RedirectUrl          string
	GithubAccessTokenAPI string
	GithubUserAPI        string
	GithubLoginUrl       string
	core.OAuthInterface
}

type GithubAuthPayload struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectUrl  string `json:"redirect_uri"`
}
