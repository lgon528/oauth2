package store

import (
	"fmt"
	"sync"

	"gopkg.in/oauth2.v3"
)

// NewUserStore create user store
func NewUserStore() *UserStore {
	return &UserStore{
		data: make(map[string]map[string]oauth2.UserInfo),
	}
}

// UserStore user information store
type UserStore struct {
	sync.RWMutex
	data map[string]map[string]oauth2.UserInfo // client_id => {user_id => userinfo}
}

// GetUser find user by clientID and userid
func (us *UserStore) GetUser(clientID string, userid string) (userinfo oauth2.UserInfo, err error) {
	us.RLock()
	defer us.RUnlock()

	if users, ok := us.data[clientID]; ok {
		if user, ok := users[userid]; ok {
			userinfo = user
			return
		}
	}
	err = fmt.Errorf("user with userid %s in clientID %s not found", userid, clientID)

	return
}

// GetUserByOpenID find user by openid
func (us *UserStore) GetUserByOpenID(clientID string, openid string) (userinfo oauth2.UserInfo, err error) {
	us.RLock()
	defer us.RUnlock()

	if users, ok := us.data[clientID]; ok {
		for _, user := range users {
			if user.GetOpenID() == openid {
				userinfo = user
				return
			}
		}
	}

	err = fmt.Errorf("user with openid %s in clientID %s not found", openid, clientID)

	return
}

// SetUser set user information
func (us *UserStore) SetUser(clientID string, userinfo oauth2.UserInfo) (err error) {
	if userinfo.GetClientID() != clientID {
		err = fmt.Errorf("clientID not match")
		return
	}

	us.Lock()
	defer us.Unlock()
	if _, ok := us.data[clientID]; !ok {
		us.data[clientID] = make(map[string]oauth2.UserInfo)
	}
	us.data[clientID][userinfo.GetID()] = userinfo

	return
}
