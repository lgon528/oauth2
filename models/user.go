package models

// User user model
type User struct {
	ID       string
	Password string
}

// GetID user id
func (u *User) GetID() string {
	return u.ID
}

// SetID user id
func (u *User) SetID(id string) {
	u.ID = id
}

// GetPassword get user password
func (u *User) GetPassword() string {
	return u.Password
}

// SetPassword set user password
func (u *User) SetPassword(pwd string) {
	u.Password = pwd
}

// OpenUser open user model
type OpenUser struct {
	User
	OpenID   string
	ClientID string
}

// GetOpenID user open id
func (u *OpenUser) GetOpenID() string {
	return u.OpenID
}

// SetOpenID user open id
func (u *OpenUser) SetOpenID(openid string) {
	u.OpenID = openid
}

// GetClientID client id
func (u *OpenUser) GetClientID() string {
	return u.ClientID
}

// SetClientID user client id
func (u *OpenUser) SetClientID(clientID string) {
	u.ClientID = clientID
}
