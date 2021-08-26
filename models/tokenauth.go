package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"neutron0.1/token"
)

type AuthToken struct {
	Id   int
	Name string
	Uuid string
	TTL  time.Duration
}

func CreateAuth(userid uint, td *token.TokenDetails) error {
	o := orm.NewOrm()
	// qs := o.QueryTable("AuthToken")

	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	auth := new(AuthToken)

	auth.Name = "AccessUuid"
	auth.Uuid = td.AccessUuid
	auth.TTL = at.Sub(now)

	_, insertErr := o.Insert(&auth)
	if insertErr != nil {
		return errors.New("failed to insert AuthToken: accessuuid")
	}

	auth.Name = "RefreshUuiD"
	auth.Uuid = td.RefreshUuid
	auth.TTL = rt.Sub(now)

	_, insertErr = o.Insert(&auth)
	if insertErr != nil {
		return errors.New("failed to insert AuthToken: refreshuuid")
	}

	return nil
}
