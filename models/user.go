package models

// User client model
type User struct {
	ID       string
	Password string
	OpenID   string
	ClientID string
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

// GetOpenID user open id
func (u *User) GetOpenID() string {
	return u.OpenID
}

// SetOpenID user open id
func (u *User) SetOpenID(openid string) {
	u.OpenID = openid
}

// GetClientID client id
func (u *User) GetClientID() string {
	return u.ClientID
}

// SetClientID user client id
func (u *User) SetClientID(clientID string) {
	u.ClientID = clientID
}
