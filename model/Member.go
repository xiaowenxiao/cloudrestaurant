package model

type Member struct {
	Id           int64   `xorm:"pk autoincr" json:"id"`
	UserName     string  `xorm:"varcher(20)" json:"user_name"`
	Mobile       string  `xorm:"varcher(11)" json:"mobile"`
	Password     string  `xorm:"varcher(255)" json:"password"`
	RegisterTime int64   `xorm:"bigint" json:"register_time"`
	Avatar       string  `xorm:"varcher(255)" json:"avatar"`
	Babance      float64 `xorm:"double" json:"balance"`
	IsActive     int8    `xorm:"tinyint" json:"is_active"`
	City         string  `xorm:"varcher(10)" json:"city"`
}
