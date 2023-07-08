package controllers

import (
	"casdoor/form"
	"casdoor/object"
	"casdoor/util"
	"fmt"
)

const SMS = "sms"

func init() {
	// 通过手机短信验证码登录
	authActionMap[SMS] = &authAction{
		login: func(authForm form.AuthForm) (user *object.User, err error) {
			if user, err = object.GetUserByFields(authForm.Organization, authForm.Username); err != nil {
				return
			} else if user == nil {
				return nil, fmt.Errorf("general:The user: %s doesn't exist", util.GetId(authForm.Organization, authForm.Username))
			}

			verificationCodeType := object.GetVerifyType(authForm.Username)
			var checkDest string

			//验证码类型==phone
			//校验号码和区号是否合法
			if verificationCodeType == object.VerifyTypePhone {
				authForm.CountryCode = user.GetCountryCode(authForm.CountryCode)
				var ok bool
				if checkDest, ok = util.GetE164Number(authForm.Username, authForm.CountryCode); !ok {
					return nil, fmt.Errorf("verification:Phone number is invalid in your region %s", authForm.CountryCode)
				}
			}

			// check result through Email or Phone
			checkResult := object.CheckSigninCode(user, checkDest, authForm.Code)
			if len(checkResult) != 0 {
				return nil, fmt.Errorf(checkResult)
			}

			// disable the verification code
			err = object.DisableVerificationCode(checkDest)
			return
		},
	}
}
