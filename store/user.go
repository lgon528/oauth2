package store

import (
	"fmt"
	"sync"

	"gopkg.in/oauth2.v3"
)

// NewMemoryUserStore create user store
func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		data: make(map[string]oauth2.UserInfo),
	}
}

// MemoryUserStore user information store
type MemoryUserStore struct {
	sync.RWMutex
	data map[string]oauth2.UserInfo // user_id => userinfo
}

// GetUser find user by userid
func (us *MemoryUserStore) GetUser(userid string) (userinfo oauth2.UserInfo, err error) {
	us.RLock()
	defer us.RUnlock()

	if user, ok := us.data[userid]; ok {
		userinfo = user
		return
	}

	err = fmt.Errorf("user %s not found", userid)

	return
}

// GetUser find user by userid
func (us *MemoryUserStore) GetUsers() []oauth2.UserInfo {
	us.RLock()
	defer us.RUnlock()

	users := make([]oauth2.UserInfo, 0)
	for _, u := range us.data {
		users = append(users, u)
	}

	return users
}

// SetUser set user information
func (us *MemoryUserStore) SetUser(userinfo oauth2.UserInfo) (err error) {
	us.Lock()
	defer us.Unlock()

	us.data[userinfo.GetID()] = userinfo

	return
}
