package domain

import validation "github.com/go-ozzo/ozzo-validation"

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r LoginReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}
