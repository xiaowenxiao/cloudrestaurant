package tool

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

type CaptchaResult struct {
	Id           string                      `json:"id"`
	DriverString *base64Captcha.DriverString // 设置生成的验证码的参数
	VertifyValue string                      `json:"code"`
	Base64Blob   string                      `json:"base_64_blob"`
}

// 默认指定本地内存存储
//var store = base64Captcha.DefaultMemStore

// 指定redis存储验证码
var store base64Captcha.Store = RedisStore{}

//生成图形化验证码
func GenerateCaptcha(ctx *gin.Context) {
	var parameters CaptchaResult = CaptchaResult{
		Id:           "",
		VertifyValue: "",
		DriverString: &base64Captcha.DriverString{
			Length:          4,
			Height:          60,
			Width:           240,
			ShowLineOptions: 2,
			NoiseCount:      0,
			Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		},
	}
	var driver base64Captcha.Driver
	driver = parameters.DriverString.ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, store)
	captchaId, base64Blob, _ := c.Generate()
	captchaResult := CaptchaResult{Id: captchaId, Base64Blob: base64Blob}

	Success(ctx, map[string]interface{}{
		"captcha_result": captchaResult,
	})

}

// base64Captcha verify
func VertifyCaptcha(id, VertifyValue string) bool {
	return store.Verify(id, VertifyValue, true)
}
