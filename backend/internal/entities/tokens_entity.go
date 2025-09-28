package entities

type AccessToken struct {
	Token string `json:"access_token" cookie:"access_token"`
}
type RefreshToken struct {
	Token string `json:"refresh_token" cookie:"refresh_token" redis:"refresh_token"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token" cookie:"access_token"`
	RefreshToken string `json:"refresh_token" cookie:"refresh_token"`
}
