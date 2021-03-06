package domain

import (
	"mime/multipart"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Apparel struct {
	ApparelID     string    `json:"apparelID"`
	ApparelTypeID string    `json:"apparelTypeID"`
	Name          string    `json:"name"`
	Brand         string    `json:"brand"`
	Price         string    `json:"price"`
	StoreURL      string    `json:"storeURL"`
	ImgURLs       []string  `json:"imageURLs"`
	CreateDate    time.Time `json:"createDate"`
	CreateBy      string    `json:"createBy"`
}

type Apparels struct {
	Apparels []*Apparel `json:"apparels"`
}

type CreateApparelReq struct {
	Requester     *Requester
	ApparelTypeID string
	Name          string
	Brand         string
	Price         string
	StoreURL      string
	Files         *multipart.Form
}

func (r CreateApparelReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ApparelTypeID, validation.Required, validation.In(ApparelTypeShirtID, ApparelTypeTShirtID, ApparelTypePantsID, ApparelTypeShortsID, ApparelTypeSkirtID, ApparelTypeDressID)),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Brand, validation.Required),
		validation.Field(&r.Price, validation.Required),
		validation.Field(&r.StoreURL, validation.Required),
		validation.Field(&r.Files, validation.Required),
	)
}

type GetApparelsAsAdminReq struct {
	Requester *Requester
	Limit     int64
	Offset    int64
}

func (r GetApparelsAsAdminReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Requester, validation.Required),
		validation.Field(&r.Limit, validation.Required, validation.Min(0), validation.Max(50)),
		validation.Field(&r.Offset, validation.Min(0)),
	)
}

type GetApparelsReq struct {
	Limit         int64
	Offset        int64
	ApparelTypeID string
}

func (r GetApparelsReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Limit, validation.Required, validation.Min(0), validation.Max(50)),
		validation.Field(&r.Offset, validation.Min(0)),
	)
}

type GetApparelByApparelIDReq struct {
	ApparelID     string
	ApparelTypeID string `json:"apparelTypeID"`
}

func (r GetApparelByApparelIDReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ApparelID, validation.Required),
	)
}

type DeleteApparelByApparelIDReq struct {
	Requester     *Requester
	ApparelID     string
	ApparelTypeID string `json:"apparelTypeID"`
}

func (r DeleteApparelByApparelIDReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Requester, validation.Required),
		validation.Field(&r.ApparelID, validation.Required),
	)
}
