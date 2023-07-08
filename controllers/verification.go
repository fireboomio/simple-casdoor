package controllers

import (
	"casdoor/form"
	"casdoor/object"
	"casdoor/util"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	LoginVerification = "login"
)

// SendVerificationCode ...
// @Title SendVerificationCode
// @Tag Verification API
// @Description "发送验证码"
// @Accept multipart/form-data
// @Param dest 				query string true "发送手机号"
// @Param countryCode 		query string false "国际区号（默认CN）" Enums(CN, US, JP) default(CN)
// @Success 200 {object} controllers.Response  "成功"
// @router /send-verification-code [post]
func (c *ApiController) SendVerificationCode() {
	var vform = form.VerificationForm{
		CaptchaType:   "none",
		Type:          "phone",
		Method:        "login",
		ApplicationId: "fireboom_fireboom_builtIn",
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &vform)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	remoteAddr := util.GetIPFromRequest(c.Ctx.Request)

	if msg := vform.CheckParameter(form.SendVerifyCode); msg != "" {
		c.ResponseError(msg)
		return
	}

	application, err := object.GetApplication(vform.ApplicationId)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	organization, err := object.GetOrganization(util.GetId(application.Owner, application.Organization))
	if err != nil {
		c.ResponseError(err.Error())
	}

	if organization == nil {
		c.ResponseError("check:Organization does not exist")
		return
	}

	var user *object.User

	sendResp := errors.New("invalid dest type")

	if vform.Method == LoginVerification {
		if user != nil && util.GetMaskedPhone(user.Phone) == vform.Dest {
			vform.Dest = user.Phone
		}

		if user, err = object.GetUserByPhone(organization.Name, vform.Dest); err != nil {
			c.ResponseError(err.Error())
			return
		} else if user == nil {
			c.ResponseError("verification:the user does not exist, please sign up first")
			return
		}

		vform.CountryCode = user.GetCountryCode(vform.CountryCode)
	}
	provider, err := object.GetProvider("fireboom/provider_sms")
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if phone, ok := util.GetE164Number(vform.Dest, vform.CountryCode); !ok {
		c.ResponseError(fmt.Sprintf("verification:Phone number is invalid in your region %s"), vform.CountryCode)
		return
	} else {
		sendResp = object.SendVerificationCodeToPhone(user, provider, remoteAddr, phone)
	}

	if sendResp != nil {
		c.ResponseError(sendResp.Error())
	} else {
		c.ResponseOk()
	}
}

func (c *ApiController) VerifyCode() {
	var authForm form.AuthForm
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &authForm)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	var user *object.User
	if authForm.Name != "" {
		user, err = object.GetUserByFields(authForm.Organization, authForm.Name)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
	}

	var checkDest string
	if strings.Contains(authForm.Username, "@") {
		if user != nil && util.GetMaskedEmail(user.Email) == authForm.Username {
			authForm.Username = user.Email
		}
		checkDest = authForm.Username
	} else {
		if user != nil && util.GetMaskedPhone(user.Phone) == authForm.Username {
			authForm.Username = user.Phone
		}
	}

	if user, err = object.GetUserByFields(authForm.Organization, authForm.Username); err != nil {
		c.ResponseError(err.Error())
		return
	} else if user == nil {
		c.ResponseError(fmt.Sprintf("general:The user: %s doesn't exist", util.GetId(authForm.Organization, authForm.Username)))
		return
	}

	verificationCodeType := object.GetVerifyType(authForm.Username)
	if verificationCodeType == object.VerifyTypePhone {
		authForm.CountryCode = user.GetCountryCode(authForm.CountryCode)
		var ok bool
		if checkDest, ok = util.GetE164Number(authForm.Username, authForm.CountryCode); !ok {
			c.ResponseError(fmt.Sprintf("verification:Phone number is invalid in your region %s"), authForm.CountryCode)
			return
		}
	}

	if result := object.CheckVerificationCode(checkDest, authForm.Code); result.Code != object.VerificationSuccess {
		c.ResponseError(result.Msg)
		return
	}
	err = object.DisableVerificationCode(checkDest)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	c.SetSession("verifiedCode", authForm.Code)

	c.ResponseOk()
}
