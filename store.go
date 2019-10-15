package oauth2

type (
	// ClientStore the client information storage interface
	ClientStore interface {
		// according to the ID for the client information
		GetByID(id string) (ClientInfo, error)

		// GetClients get all clients
		GetClients() []ClientInfo

		// AddClient add a client
		AddClient(info ClientInfo) error
	}

	// TokenStore the token information storage interface
	TokenStore interface {
		// create and store the new token information
		Create(info TokenInfo) error

		// delete the authorization code
		RemoveByCode(code string) error

		// use the access token to delete the token information
		RemoveByAccess(access string) error

		// use the refresh token to delete the token information
		RemoveByRefresh(refresh string) error

		// use the authorization code for token information data
		GetByCode(code string) (TokenInfo, error)

		// use the access token for token information data
		GetByAccess(access string) (TokenInfo, error)

		// use the refresh token for token information data
		GetByRefresh(refresh string) (TokenInfo, error)
	}

	// UserStore the user information storage interface
	UserStore interface {
		// GetUser find user by clientID and userid
		GetUser(userid string) (UserInfo, error)

		// GetUsers get all user
		GetUsers() []UserInfo

		// SetUser set user information
		SetUser(userinfo UserInfo) error
	}
)
