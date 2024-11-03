package core

import (
	"errors"
	"log"
)

func (o *OAuth) AddProvider(name string, provider OAuthInterface) {
	if o.providers == nil {
		o.providers = make(map[string]OAuthInterface)
	}
	if o.providers[name] != nil {
		panic("already using name")
	}
	o.providers[name] = provider
}

func (o *OAuth) GetProviders() ([]Provider, error) {
	if o.providers == nil {
		o.providers = make(map[string]OAuthInterface)
	}

	providers := []Provider{}

	for name, provider := range o.providers {
		providers = append(providers, Provider{Name: name, LoginUrl: provider.GetLoginUrl()})
	}

	return providers, nil
}

func (o *OAuth) GetLoginUrl(name string) (string, error) {
	if o.providers == nil {
		o.providers = make(map[string]OAuthInterface)
	}

	if o.providers[name] == nil {
		log.Fatalln("no such provider")
		return "", errors.New("no such provider")
	}

	return o.providers[name].GetLoginUrl(), nil
}

func (o *OAuth) GetUser(name string, code string) (*User, error) {
	if o.providers == nil {
		o.providers = make(map[string]OAuthInterface)
	}

	if o.providers[name] == nil {
		log.Fatalln("no such provider")
		return nil, errors.New("no such provider")
	}

	return o.providers[name].GetUser(code)
}

var Core = OAuth{}
