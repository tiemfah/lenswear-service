package uidgen

import "github.com/golang-plus/uuid"

type UIDGen interface {
	New(key string) string
}

type uidgen struct{}

func New() UIDGen {
	return &uidgen{}
}

func (u uidgen) New(key string) string {
	uid, _ := uuid.NewV4()
	return key + uid.Format(uuid.StyleStandard)
}
