package cache

import (
	"me.qqtu.game/model"
)

type UserData struct {
	UserId  int64               //用户id
	account *model.AccountModel //账号数据
}

func (ud *UserData) GetAccountModel() *model.AccountModel {
	if ud.account == nil {
		ud.account = &model.AccountModel{UserId: ud.UserId}
		ud.account.GetByUserId()
	}
	return ud.account
}
