package common

import (
	"errors"

	"github.com/jinzhu/gorm"

	base "github.com/shiniu0606/engine/core/base"
	jbp "github.com/shiniu0606/engine/core/db"
)

func CreateServer(server *Server) (int64,error) {
	//校验参数
	if len(server.ServerAddr) == 0 || server.ServerId == 0 {
		return 0, errors.New("CreateServer param error")
	}

	if err := jbp.GetDB().First(server, "server_id=? ", server.ServerId).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return 0, errors.New("CreateServer server_id query error")
		}
	} else {
		//已存在直接返回
		return server.UID, nil
	}

	if err := jbp.GetDB().Create(server).Error; err != nil {
		base.LogError("CreateServer data failed:%v",err)
		return 0, errors.New("CreateServer data failed")
	}
	return server.UID, nil
}

func UpdateOnlineByUid(Uid int64,num int) (error) {
	data := &Server{}
	if err := jbp.GetDB().First(&data, "uid=? ", Uid).Error; err != nil {
		return errors.New("UpdateOnlineByUid error:" + err.Error())
	}

	if data.UID < 1 {
		return errors.New("UpdateOnlineByUid error")
	}

	err := jbp.GetDB().Model(&data).Update("online_num", num).Error
	if err != nil {
		return errors.New("UpdateOnlineByUid failed:" + err.Error())
	}

	return nil
}