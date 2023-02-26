package controller

import (
	"cloudrestaurant/param"
	"cloudrestaurant/service"
	"cloudrestaurant/tool"
	"fmt"

	"github.com/gin-gonic/gin"
)

type MemberController struct {
}

func (mc *MemberController) Router(engine *gin.Engine) {
	engine.GET("/api/sendcode", mc.sendSmsCode)
	engine.OPTIONS("/api/login", mc.smsLogin)
	engine.GET("/api/captcha", mc.captcha)
}

// 生成验证码
func (mc *MemberController) captcha(context *gin.Context) {
	// todo 生成验证码，并返回客户端
	tool.GenerateCaptcha(context)
}

// 判断验证码是否正确
func (ms *MemberController) vertifyCaptcha(context *gin.Context) {
	var captcha tool.CaptchaResult
	err := tool.Decode(context.Request.Body, &captcha)
	if err != nil {
		tool.Failed(context, "参数解析失败")
		return
	}
	result := tool.VertifyCaptcha(captcha.Id, captcha.VertifyValue)
	if result {
		fmt.Println("验证通过")
	} else {
		fmt.Println("验证失败")
	}

}

func (mc *MemberController) sendSmsCode(context *gin.Context) {
	// 发送验证码
	phone, exist := context.GetQuery("phone")
	if !exist {
		tool.Failed(context, "参数解析失败")
		return
	}
	ms := service.MemberService{}
	isSend := ms.Sendcode(phone)
	if isSend {
		tool.Success(context, "发送成功")
		return
	}
	tool.Failed(context, "发送失败")
}

// 手机号+短信 登录方法
func (mc *MemberController) smsLogin(context *gin.Context) {
	var SmsLoginParam param.SmsLoginParam
	err := tool.Decode(context.Request.Body, &SmsLoginParam)
	if err != nil {
		tool.Failed(context, "参数解析失败")
		return
	}
	// 完成手机+验证码登录
	us := service.MemberService{}
	member := us.SmsLogin(SmsLoginParam)
	if member != nil {
		tool.Success(context, member)
		return
	}
	tool.Failed(context, "登录失败")
}
