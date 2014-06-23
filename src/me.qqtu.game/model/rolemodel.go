package model

import (
	"labix.org/v2/mgo/bson"
	"me.qqtu.game/mio"
)

type RoleModel struct {
	UserId int64
	Sex    int
}

func (r *RoleModel) Add() error {
	return mio.GetRoleC().Insert(r)
}

func (r *RoleModel) Update() error {
	return mio.GetRoleC().Update(bson.M{"userid": r.UserId}, r)
}

func (r *RoleModel) Get() error {
	return mio.GetRoleC().Find(bson.M{"userid": r.UserId}).One(r)
}
