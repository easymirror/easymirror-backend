package cyberfile

import "context"

type Account interface {
	Username() string
	Password() string
	AccessToken() string
	AccountID() string
	GetAccessToken(ctx context.Context) (string, error)
}

type account struct {
	username, password string
	accountID          string
	accessToken        string
}

// NewAccount is a wrapper function to create a new account
func NewAccount(username, password string) Account {
	return newAccount(username, password)
}

// newAccount creates a new account with given credentials
func newAccount(username, password string) Account {
	return &account{
		username: username,
		password: password,
	}
}

// Username is a getter function that returns the username
func (a account) Username() string { return a.username }

// Password is a getter function that returns the password
func (a account) Password() string { return a.password }

// AccessToken is a getter function that returns the access token
func (a account) AccessToken() string { return a.accessToken }

// AccountID is a getter function that returns the access token
func (a account) AccountID() string { return a.accountID }
