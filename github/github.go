package github

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/schedule-job/schedule-job-authorization/core"
	schedule_errors "github.com/schedule-job/schedule-job-errors"
)

func (g *Github) getAccessToken(code string) (string, error) {
	payload := GithubAuthPayload{
		ClientId:     g.ClientId,
		ClientSecret: g.ClientSecret,
		Code:         code,
		RedirectUrl:  g.RedirectUrl,
	}

	body, errMarshal := json.Marshal(payload)

	if errMarshal != nil {
		err := schedule_errors.InvalidArgumentError{Param: "payload", Message: errMarshal.Error()}
		log.Fatalln(err.Error())
		return "", &err
	}

	req, errReq := http.NewRequest("POST", g.GithubAccessTokenAPI, bytes.NewReader(body))

	if errReq != nil {
		err := schedule_errors.ConnectionError{Address: g.GithubAccessTokenAPI, Reason: errReq.Error()}
		log.Fatalln(err.Error())
		return "", &err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, errRes := client.Do(req)

	if errRes != nil {
		err := schedule_errors.ConnectionError{Address: g.GithubAccessTokenAPI, Reason: errRes.Error()}
		log.Fatalln(err.Error())
		return "", &err
	}

	defer res.Body.Close()

	var userData map[string]interface{}
	errDecode := json.NewDecoder(res.Body).Decode(&userData)

	if errDecode != nil {
		err := schedule_errors.InternalServerError{Err: errDecode}
		log.Fatalln(err.Error())
		return "", &err
	}

	if userData["error"] != "" && userData["error"] != nil {
		err := schedule_errors.UnauthorizedError{Reason: userData["error_description"].(string) + " more info : " + userData["error_uri"].(string)}
		log.Fatalln(err)
		return "", &err
	}

	return userData["access_token"].(string), nil
}

func (g *Github) getUser(accessToken string) (*core.User, error) {
	req, errReq := http.NewRequest("GET", g.GithubUserAPI, nil)

	if errReq != nil {
		err := schedule_errors.ConnectionError{Address: g.GithubAccessTokenAPI, Reason: errReq.Error()}
		log.Fatalln(err.Error())
		return nil, &err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	res, errRes := client.Do(req)

	if errRes != nil {
		err := schedule_errors.ConnectionError{Address: g.GithubAccessTokenAPI, Reason: errRes.Error()}
		log.Fatalln(err.Error())
		return nil, &err
	}

	read, errRead := io.ReadAll(res.Body)

	if errRead != nil {
		err := schedule_errors.InternalServerError{Err: errRead}
		log.Fatalln(err.Error())
		return nil, &err
	}

	user := core.User{}

	errParse := json.Unmarshal(read, &user)

	if errParse != nil {
		err := schedule_errors.InternalServerError{Err: errParse}
		log.Fatalln(err.Error())
		return nil, &err
	}

	return &user, nil
}

func (g *Github) GetUser(code string) (*core.User, error) {
	accessToken, errAccessToken := g.getAccessToken(code)

	if errAccessToken != nil {
		return nil, errAccessToken
	}

	return g.getUser(accessToken)
}

func (g *Github) GetLoginUrl() string {
	return g.GithubLoginUrl + "?scope=user&client_id=" + g.ClientId + "&redirect_url=" + g.RedirectUrl
}
