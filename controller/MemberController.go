package controller

import (
	"cloudrestaurant/model"
	"cloudrestaurant/param"
	"cloudrestaurant/service"
	"cloudrestaurant/tool"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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
	// 头像上传
	engine.POST("/api/upload/avator", mc.uploadAvator)
}

func (mc *MemberController) uploadAvator(context *gin.Context) {
	// 1、解析上传的参数：file、user_id
	userId := context.PostForm("user_id")
	fmt.Println(userId)
	file, err := context.FormFile("avatar")
	if err != nil || userId == "" {
		tool.Failed(context, "参数解析失败")
		return
	}

	// 2、判断user_id对应的用户是否已经登录
	sess := tool.Getsess(context, "user_"+userId)
	if sess == nil {
		tool.Failed(context, "参数不合法")
		return
	}
	var member model.Member
	json.Unmarshal(sess.([]byte), &member)

	// 3、file保存到本地
	fileName := "./uploadfile/" + strconv.FormatInt(time.Now().Unix(), 10) + file.Filename
	err = context.SaveUploadedFile(file, fileName)
	if err != nil {
		tool.Failed(context, "头像更新失败")
		return
	}

	// 4、将保存后的文件本地路径，保存到用户表中的头像字段
	memberService := service.MemberService{}
	path := memberService.UploadAvator(member.Id, fileName[1:])
	if path != "" {
		tool.Success(context, "http://localhost:8090"+path)
		return
	}

	// 5、返回结果
	return
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
		// 用户信息保存到session
		sess, _ := json.Marshal(member)
		err = tool.Setsess(context, "user_"+string(rune(member.Id)), sess)
		if err != nil {
			tool.Failed(context, "登录失败")
			return
		}
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
		// 用户信息保存到session
		sess, _ := json.Marshal(member)
		err = tool.Setsess(context, "user_"+string(rune(member.Id)), sess)
		if err != nil {
			tool.Failed(context, "登录失败")
			return
		}
		tool.Success(context, member)
		return
	}
	tool.Failed(context, "登录失败")
}
