package auth

const (
	RefreshCookieName = "jwt_refresh"
)

type AuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
