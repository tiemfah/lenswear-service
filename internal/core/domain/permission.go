package domain

import validation "github.com/go-ozzo/ozzo-validation"

type Permission struct {
	PermissionID  string `json:"permissionID"`
	TransformerID string `json:"transformerID"`
}

type Permissions struct {
	Permissions []*Permission
}

type CreatePermissionReq struct {
	Requester     *Requester
	Name          string `json:"name"`
	TransformerID string `json:"transformerID"`
}

func (r CreatePermissionReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Requester, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.TransformerID, validation.Required),
	)
}

type GetPermissionsAsAdminReq struct {
	Requester *Requester
	Limit     int64
	Offset    int64
}

func (r GetPermissionsAsAdminReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Requester, validation.Required),
		validation.Field(&r.Limit, validation.Required, validation.Min(0), validation.Max(50)),
		validation.Field(&r.Offset, validation.Min(0)),
	)
}

type GetPermissionsByUserIDReq struct {
	Requester *Requester
}

func (r GetPermissionsByUserIDReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Requester, validation.Required),
	)
}

type AddPermissionToUserReq struct {
	Requester    *Requester
	PermissionID string `json:"permissionID"`
	UserID       string `json:"userID"`
}

func (r AddPermissionToUserReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Requester, validation.Required),
		validation.Field(&r.PermissionID, validation.Required),
		validation.Field(&r.UserID, validation.Required),
	)
}

type RemovePermissionFromUserReq struct {
	Requester    *Requester
	PermissionID string `json:"permissionID"`
	UserID       string `json:"userID"`
}

func (r RemovePermissionFromUserReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Requester, validation.Required),
		validation.Field(&r.PermissionID, validation.Required),
		validation.Field(&r.UserID, validation.Required),
	)
}

type DeletePermissionReq struct {
	Requester    *Requester
	PermissionID string `json:"permissionID"`
}

func (r DeletePermissionReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.PermissionID, validation.Required),
	)
}

// helpers

type CheckUserIsPermittedOnTransformerReq struct {
	UserID        string
	TransformerID string
}
