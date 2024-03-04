package cyberfile

type Account interface {
	Username() string
	Password() string
	AccountID() string
	GetAccessToken() (string, error)
}

type account struct {
	username, password string
	accountID          string
	accessToken        string
}
