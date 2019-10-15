package models

import (
	"gopkg.in/oauth2.v3"
	"sync"
)

// Client client model
type Client struct {
	ID          string
	Secret      string
	Domain      string
	Name        string
	Description string

	mutex sync.RWMutex
	Users map[string]oauth2.OpenUserInfo
}

func (c *Client) SetID(id string) {
	c.ID = id
}

func (c *Client) SetSecret(secret string) {
	c.Secret = secret
}

func (c *Client) GetName() string {
	return c.Name
}

// GetID client id
func (c *Client) GetID() string {
	return c.ID
}

// GetSecret client domain
func (c *Client) GetSecret() string {
	return c.Secret
}

// GetDomain client domain
func (c *Client) GetDomain() string {
	return c.Domain
}

// GetOpenUser get user info
func (c *Client) GetOpenUser(userid string) (user oauth2.OpenUserInfo, ok bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.Users == nil {
		c.Users = make(map[string]oauth2.OpenUserInfo)
		ok = false
		return
	}

	user, ok = c.Users[userid]

	return
}

func (c *Client) SetOpenUser(user oauth2.OpenUserInfo) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.Users == nil {
		c.Users = make(map[string]oauth2.OpenUserInfo)
	}

	c.Users[user.(oauth2.UserInfo).GetID()] = user
	return
}
