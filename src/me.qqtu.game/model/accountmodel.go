package model

import (
	"labix.org/v2/mgo/bson"
	"me.qqtu.game/logger"
	"me.qqtu.game/mio"
	"me.qqtu.game/rio"
	"me.qqtu.game/router"
)

type AccountModel struct {
	UserId   int64
	Account  string
	Password string
	Online   bool
}

func (a *AccountModel) Add() error {
	return mio.GetAccountC().Insert(a)
}

func (a *AccountModel) Update() error {
	return mio.GetAccountC().Update(bson.M{"account": a.Account, "userid": a.UserId}, a)
}

func (a *AccountModel) IsExists() (bool, error) {
	count, err := mio.GetAccountC().Find(bson.M{"$or": []bson.M{bson.M{"account": a.Account, "userid": a.UserId}}}).Count()
	return count > 0, err
}

func (a *AccountModel) GetByAccount() bool {
	if err := mio.GetAccountC().Find(bson.M{"account": a.Account}).One(a); err != nil {
		return false
	}
	return true
}

func (a *AccountModel) GetByUserId() bool {
	if err := mio.GetAccountC().Find(bson.M{"userid": a.UserId}).One(a); err != nil {
		return false
	}
	return true
}

func (a *AccountModel) Login() bool {
	acc := a.Account
	pwd := a.Password

	ret, err := a.IsExists()
	if err != nil {
		logger.GetLogger().Error(err.Error(), nil)
		return false
	}
	if ret {
		a.GetByAccount()
	}
	if a.UserId > 0 && a.Password == pwd && a.Account == acc {
		a.Online = true
		a.Update()
		return true
	}
	return false
}

func (a *AccountModel) Register() bool {
	var err error
	var isExists bool
	var dbRedis *rio.DBRedis

	isExists, err = a.IsExists()

	if isExists || err != nil {
		return false
	}

	dbRedis = rio.GetReidsByServerId(router.CurrSID)

	if dbRedis == nil {
		return false
	}
	a.UserId, err = dbRedis.GetNextUserId()
	if err != nil {
		return false
	}

	err = a.Add()

	return err == nil
}

func (a *AccountModel) Logout() bool {
	if a.GetByUserId() {
		a.Online = false
		err := a.Update()
		return err == nil
	}
	return false
}
