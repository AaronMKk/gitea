package user

import (
	"errors"
	"github.com/opensourceways/xihe-server/utils"
	"strings"
)

// Account
type Account interface {
	Account() string
}

type dpAccount string

type oldUserTokenPayload struct {
	Account                 string `json:"account"`
	Email                   string `json:"email"`
	PlatformToken           string `json:"token"`
	PlatformUserNamespaceId string `json:"nid"`
}

func (r dpAccount) Account() string {
	return string(r)
}

func (pl *oldUserTokenPayload) hasEmail() bool {
	return pl.Email != "" && pl.PlatformToken != ""
}

func NewAccount(v string) (Account, error) {
	if v == "" || strings.ToLower(v) == "root" || !utils.IsUserName(v) {
		return nil, errors.New("invalid user name")
	}

	return dpAccount(v), nil
}

func (pl *oldUserTokenPayload) DomainAccount() Account {
	a, _ := NewAccount(pl.Account)

	return a
}
