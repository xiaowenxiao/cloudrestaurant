package model

type SmsCode struct {
	Id         int    `xorm:"pk autoincr" json:"id"`
	Phone      string `xorm:"varchar(11)" json:"phone"`
	BizId      string `xorm:"varchar(30)" json:"bizid"`
	Code       string `xorm:"varchar(6)" json:"code"`
	CreateTime int64  `xorm:"bigint" json:"create_time"`
}
