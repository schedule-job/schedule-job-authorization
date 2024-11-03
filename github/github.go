package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/schedule-job/schedule-job-authorization/core"
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
		log.Fatalln(errMarshal.Error())
		return "", errMarshal
	}

	req, errReq := http.NewRequest("POST", g.GithubAccessTokenAPI, bytes.NewReader(body))

	if errReq != nil {
		log.Fatalln(errReq.Error())
		return "", errReq
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, errRes := client.Do(req)

	if errRes != nil {
		log.Fatalln(errRes.Error())
		return "", errRes
	}

	defer res.Body.Close()

	var userData map[string]interface{}
	errDecode := json.NewDecoder(res.Body).Decode(&userData)

	if errDecode != nil {
		log.Fatalln(errDecode.Error())
		return "", errDecode
	}

	if userData["error"] != "" && userData["error"] != nil {
		log.Fatalln(userData["error_description"].(string) + " more info : " + userData["error_uri"].(string))
		return "", errors.New(userData["error_description"].(string) + " more info : " + userData["error_uri"].(string))
	}

	return userData["access_token"].(string), nil
}

func (g *Github) getUser(accessToken string) (*core.User, error) {
	req, errReq := http.NewRequest("GET", g.GithubUserAPI, nil)

	if errReq != nil {
		log.Fatalln(errReq.Error())
		return nil, errReq
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	res, errRes := client.Do(req)

	if errRes != nil {
		log.Fatalln(errRes.Error())
		return nil, errRes
	}

	read, errRead := io.ReadAll(res.Body)

	if errRead != nil {
		log.Fatalln(errRead.Error())
		return nil, errRead
	}

	user := core.User{}

	errParse := json.Unmarshal(read, &user)

	if errParse != nil {
		log.Fatalln(errParse.Error())
		return nil, errParse
	}

	return &user, nil
}

func (g *Github) GetUser(code string) (*core.User, error) {
	accessToken, errAccessToken := g.getAccessToken(code)

	if errAccessToken != nil {
		log.Fatalln(errAccessToken.Error())
		return nil, errAccessToken
	}

	return g.getUser(accessToken)
}

func (g *Github) GetLoginUrl() string {
	return g.GithubLoginUrl + "?scope=user&client_id=" + g.ClientId + "&redirect_url=" + g.RedirectUrl
}
