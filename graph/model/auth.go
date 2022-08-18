package model

type Credential struct {
	AccessToken   string `json:"accessToken"`
	RefreshToken  string `json:"refreshToken"`
	AccessExpire  int64  `json:"accessExpire"`
	RefreshExpire int64  `json:"refreshExpire"`
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
