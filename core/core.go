package core

import (
	"log"

	schedule_errors "github.com/schedule-job/schedule-job-errors"
)

func (o *OAuth) AddProvider(name string, provider OAuthInterface) error {
	if o.providers == nil {
		o.providers = make(map[string]OAuthInterface)
	}

	if o.providers[name] != nil {
		err := schedule_errors.DuplicateNameError{Name: name}
		log.Fatalln(err.Error())
		return &err
	}

	o.providers[name] = provider

	return nil
}

func (o *OAuth) GetProviders() []Provider {
	if o.providers == nil {
		o.providers = make(map[string]OAuthInterface)
	}

	providers := []Provider{}

	for name, provider := range o.providers {
		providers = append(providers, Provider{Name: name, LoginUrl: provider.GetLoginUrl()})
	}

	return providers
}

func (o *OAuth) GetLoginUrl(name string) (string, error) {
	if o.providers == nil {
		o.providers = make(map[string]OAuthInterface)
	}

	if o.providers[name] == nil {
		err := schedule_errors.UnsupportedFeatureError{Feature: name}
		log.Fatalln(err.Error())
		return "", &err
	}

	return o.providers[name].GetLoginUrl(), nil
}

func (o *OAuth) GetUser(name string, code string) (*User, error) {
	if o.providers == nil {
		o.providers = make(map[string]OAuthInterface)
	}

	if o.providers[name] == nil {
		err := schedule_errors.UnsupportedFeatureError{Feature: name}
		log.Fatalln(err.Error())
		return nil, &err
	}

	return o.providers[name].GetUser(code)
}

var Core = OAuth{}
