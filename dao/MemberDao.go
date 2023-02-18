package dao

import (
	"cloudrestaurant/model"
	"cloudrestaurant/tool"
	"fmt"

	"github.com/wonderivan/logger"
)

type MemberDao struct {
	*tool.Orm
}

// 验证手机号和验证码是否存在
func (md *MemberDao) ValiddateSmsCode(phone string, code string) *model.SmsCode {
	var sms model.SmsCode

	if _, err := md.Where("phone = ? and code = ? ", phone, code).Get(&sms); err != nil {
		fmt.Println(err.Error())
	}
	return &sms
}

// 查询member表中phone是否存在
func (md *MemberDao) QueryByPhone(phone string) *model.Member {
	var member model.Member
	if _, err := md.Where("mobile = ? ", phone).Get(&member); err != nil {
		fmt.Println(err.Error())
	}
	return &member
}

// 新用户的数据库插入操作
func (md *MemberDao) InsertMember(member model.Member) int64 {
	result, err := md.InsertOne(&member)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	return result
}

// 将验证码插入数据库
func (md *MemberDao) InsertCode(sms model.SmsCode) int64 {
	result, err := md.InsertOne(&sms)
	if err != nil {
		logger.Error(err.Error())
	}
	return result
}
