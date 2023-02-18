package dao

import (
	"cloudrestaurant/model"
	"cloudrestaurant/tool"

	"github.com/wonderivan/logger"
)

type MemberDao struct {
	*tool.Orm
}

func (md *MemberDao) InsertCode(sms model.SmsCode) int64 {
	result, err := md.InsertOne(&sms)
	if err != nil {
		logger.Error(err.Error())
	}
	return result
}
