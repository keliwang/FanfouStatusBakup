package main

import (
	"encoding/json"
	"errors"
	"github.com/mrjones/oauth"
	"io/ioutil"
	"os"
)

// Fanfou's OAuth authention urls
var fanfouService = &oauth.ServiceProvider{
	RequestTokenUrl:   "http://fanfou.com/oauth/request_token",
	AuthorizeTokenUrl: "http://fanfou.com/oauth/authorize",
	AccessTokenUrl:    "http://fanfou.com/oauth/access_token",
}

// A struct to store Fanfou's OAuth information
type Fanfou struct {
	Key         string
	Secret      string
	AccessToken *oauth.AccessToken
}

// Fanfou Client struct
type FanfouClient struct {
	*oauth.Consumer
	Token *oauth.AccessToken
}

// Check if we have our access token already
func (c *Fanfou) CheckAccessToken() bool {
	if c.AccessToken == nil {
		return false
	}

	// We have access token already
	return true
}

// Load config from a json file
func (c *Fanfou) LoadConfig(name string) error {
	jsonBytes, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, c)
	if err != nil {
		return err
	}

	return nil
}

// Write config to a json file
func (c *Fanfou) WriteConfig(name string) error {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(name, jsonBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Create a new Fanfou client, root should be RootFanfou or RootSandbox
func (c *Fanfou) NewFanfouClient() (*FanfouClient, error) {
	if c.Key == "" || c.Secret == "" {
		return nil, errors.New("No consumer key or secret")
	}

	return &FanfouClient{
		oauth.NewConsumer(c.Key, c.Secret, *fanfouService),
		c.AccessToken,
	}, nil
}

// A wrapper for oauth.Consumer.AuthorizeToken. Please remember to store your
// access token to Fanfou's OAuth information struct after call this function.
func (c *FanfouClient) AuthorizeToken(
	rtoken *oauth.RequestToken,
	verifyCode string) (*oauth.AccessToken, error) {

	atoken, err := c.Consumer.AuthorizeToken(rtoken, verifyCode)
	if err == nil {
		// store access token
		c.Token = atoken
	}

	return atoken, err
}
