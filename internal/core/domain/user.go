package domain

import validation "github.com/go-ozzo/ozzo-validation"

type User struct {
	UserID     string `json:"userID"`
	UserRoleID string `json:"userRoleID"`
	Username   string `json:"username"`
	AreaID     string `json:"areaID"`
}

type Users struct {
	Users []*User `json:"users"`
}

type CreateUserReq struct {
	Requester *Requester
	Username  string `json:"username"`
	Password  string `json:"password"`
	AreaID    string `json:"areaID"`
}

func (r CreateUserReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Password, validation.Required),
		validation.Field(&r.AreaID, validation.Required),
	)
}

type GetUsersAsAdminReq struct {
	Requester *Requester
	Limit     int64
	Offset    int64
}

func (r GetUsersAsAdminReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Requester, validation.Required),
		validation.Field(&r.Limit, validation.Required, validation.Min(0), validation.Max(50)),
		validation.Field(&r.Offset, validation.Min(0)),
	)
}

type GetUserByUserIDReq struct {
	UserID string `json:"userID"`
}

type ModifyUserReq struct {
	Requester *Requester
	Username  string `json:"username"`
	AreaID    string `json:"areaID"`
}

func (r ModifyUserReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Requester, validation.Required),
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.AreaID, validation.Required),
	)
}

type ResetUserPasswordReq struct {
	Requester *Requester
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (r ResetUserPasswordReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Requester, validation.Required),
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}

type CheckUserIsValidPasswordReq struct {
	Username string
	Password string
}
