package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"geek-nebula/libs"
	"geek-nebula/models"
	"geek-nebula/utils"
	"github.com/astaxie/beego"
	beegoCache "github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/utils/captcha"

	cache "github.com/patrickmn/go-cache"
)

type LoginController struct {
	BaseController
}

var cpt *captcha.Captcha

func init() {
	// use beego cache system store the captcha data
	store := beegoCache.NewMemoryCache()
	cpt = captcha.NewWithFilter("/captcha/", store)
	cpt.ChallengeNums = 5
	cpt.StdWidth = 180
	cpt.StdHeight = 40
	cpt.FieldCaptchaName = "authCode"
	cpt.FieldIDName = "captcha_id"
}

//登录 TODO:XSRF过滤
func (self *LoginController) LoginIn() {
	if self.userId > 0 {
		self.redirect(beego.URLFor("HomeController.Index"))
	}
	beego.ReadFromRequest(&self.Controller)
	if self.isPost() {

		username := strings.TrimSpace(self.GetString("username"))
		password := strings.TrimSpace(self.GetString("password"))
		//authCode := strings.TrimSpace(self.GetString("authCode"))
		//captchaId := strings.TrimSpace(self.GetString("captcha_id"))

		if username != "" && password != "" {
			user, err := models.AdminGetByName(username)
			fmt.Println(user)
			flash := beego.NewFlash()
			errorMsg := ""
			if err != nil || user.Password != libs.Md5([]byte(password+user.Salt)) {
				errorMsg = "帐号或密码错误"
			} else if user.Status == 0 {
				errorMsg = "该帐号已禁用"
			} else if !cpt.VerifyReq(self.Ctx.Request) {
				//!cpt.Verify(captchaId, authCode)
				errorMsg = "验证码错误"
			} else {
				user.LastIp = self.getClientIp()
				user.LastLogin = time.Now().Unix()
				user.Update()
				utils.Che.Set("uid"+strconv.Itoa(user.Id), user, cache.DefaultExpiration)
				authkey := libs.Md5([]byte(self.getClientIp() + "|" + user.Password + user.Salt))
				self.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 7*86400)

				self.redirect(beego.URLFor("HomeController.Index"))
			}
			flash.Error(errorMsg)
			flash.Store(&self.Controller)
			self.redirect(beego.URLFor("LoginController.LoginIn"))
		}
	}
	self.TplName = "login/login.html"
}

//登出
func (self *LoginController) LoginOut() {
	self.Ctx.SetCookie("auth", "")
	self.redirect(beego.URLFor("LoginController.LoginIn"))
}

func (self *LoginController) NoAuth() {
	self.Ctx.WriteString("没有权限")
}
