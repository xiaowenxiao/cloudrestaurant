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
	// 发送短信验证码，调用阿里云sdk发送
	engine.GET("/api/sendcode", mc.sendSmsCode)
	// 手机号+短信验证码登录
	engine.POST("/api/login_sms", mc.smsLogin)

	// 生成base64Captcha验证码
	engine.GET("/api/captcha", mc.captcha)
	// 验证base64Captcha验证码 postman测试
	engine.POST("/api/vertifycha", mc.vertifyCaptcha)
	// 账号密码登录
	engine.POST("/api/login_pwd", mc.nameLogin)
}

func (mc *MemberController) nameLogin(context *gin.Context) {
	// 1、解析用户登录传递参数
	var loginparam param.LoginParam
	err := tool.Decode(context.Request.Body, &loginparam)
	if err != nil {
		tool.Failed(context, "参数解析失败")
		return
	}

	// 2、验证验证码
	validate := tool.VertifyCaptcha(loginparam.Id, loginparam.Value)
	if !validate {
		tool.Failed(context, "验证码不正确，请重新验证")
		return
	}

	// 3、登录
	ms := service.MemberService{}
	member := ms.Login(loginparam.Name, loginparam.Password)
	if member.Id != 0 {
		tool.Success(context, &member)
		return
	}
	tool.Failed(context, "登录失败")
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
