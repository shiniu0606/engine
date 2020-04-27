package common

import (
	"errors"

	"github.com/jinzhu/gorm"

	base "github.com/shiniu0606/engine/core/base"
	jbp "github.com/shiniu0606/engine/core/db"
)

func CreateAccount(account *Account) (int64,error) {
	//校验参数
	if len(account.AccountName)  == 0 || account.UserId == 0  {
		return 0, errors.New("CreateAccount param error")
	}

	if err := jbp.GetDB().Create(account).Error; err != nil {
		base.LogError("CreateAccount data failed:%v",err)
		return 0, errors.New("CreateAccount data failed")
	}
	return account.UID, nil
}

func GetAccountByUserId(userid int64) (*Account,error) {
	account := Account{}
	if err := jbp.GetDB().First(&account, "user_id=? ", userid).Error; err != nil {
		return nil, errors.New("CreateAccount user_id query error")
	} 
	return &account, nil
}

func GetAccountByAccountName(userid int64) (*Account,error) {
	account := Account{}
	if err := jbp.GetDB().First(account, "acc_name=? ", account.AccountName).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, errors.New("CreateAccount user_id query error")
		}
		return nil,nil
	}
	return &account, nil 
}
