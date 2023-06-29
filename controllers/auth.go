package controllers

import (
	"casdoor/form"
	"casdoor/object"
	"casdoor/util"
	"encoding/json"
	"fmt"
	"sync"
)

var (
	lock sync.RWMutex
)

func codeToResponse(code *object.Code) *Response {
	if code.Code == "" {
		return &Response{Status: "error", Msg: code.Message, Data: code.Code}
	}

	return &Response{Status: "ok", Msg: "", Data: code.Code}
}

// 返回 AccessToken 和 RefreshToken
type TokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func tokenToResponse(token *object.Token) *UserResponse {
	if token.AccessToken == "" {
		return &UserResponse{Success: false, Data: TokenResp{}}
	}
	return &UserResponse{Success: true, Data: TokenResp{token.AccessToken, token.RefreshToken}}
}

func (c *ApiController) HandleLoggedIn(application *object.Application, user *object.User, form *form.AuthForm) (resp *UserResponse) {
	if form.Type == ResponseTypeToken || form.Type == ResponseTypeIdToken {
		token, _ := object.GetTokenByUser(application, user, "", c.Ctx.Request.Host)
		resp = tokenToResponse(token)
	} else {
		resp = wrapErrorUserResponse(fmt.Errorf("unknown response type: %s", form.Type))
	}
	// if user did not check auto signin
	if resp.Success && !form.AutoSignin {
		c.setExpireForSession()
	}
	if resp.Success {
		_, err := object.AddSession(&object.Session{
			Owner:       user.Owner,
			Name:        user.Name,
			Application: application.Name,
			SessionId:   []string{c.Ctx.Input.CruSession.SessionID()},
		})
		if err != nil {
			c.ResponseError(err.Error(), nil)
			return
		}
	}
	return resp
}

// Login ...
// @Title Login
// @Tag Login API
// @Description login
// @Param username        query    string  true "用户名/号码"
// @Param organization    query    string  true "组织"
// @Param countryCode     query    string  false "国际区号（默认CN）" Enums(CN, US, JP) default(CN)
// @Param code     		  query    string  true  "验证码"
// @Param type     		  query    string  true  "类型：token"
// @Param application     query    string  true  "应用名称"
// @Success 200 {object} controllers.UserResponse  		成功
// @router /login [post]
func (c *ApiController) Login() {
	resp := &UserResponse{}

	var authForm form.AuthForm
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &authForm)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if authForm.Username != "" {
		if authForm.Type == ResponseTypeLogin {
			if c.GetSessionUsername() != "" {
				c.Response(false, "account:Please sign out first", TokenResp{})
				return
			}
		}

		var user *object.User

		// 密码为空--验证码登录
		if authForm.Password == "" {
			if user, err = object.GetUserByFields(authForm.Organization, authForm.Username); err != nil {
				c.Response(false, err.Error(), TokenResp{})
				return
			} else if user == nil {
				c.Response(false, fmt.Sprintf("general:The user: %s doesn't exist", util.GetId(authForm.Organization, authForm.Username)), TokenResp{})
				return
			}

			verificationCodeType := object.GetVerifyType(authForm.Username)
			var checkDest string

			//验证码类型==phone
			//校验号码和区号是否合法
			if verificationCodeType == object.VerifyTypePhone {
				authForm.CountryCode = user.GetCountryCode(authForm.CountryCode)
				var ok bool
				if checkDest, ok = util.GetE164Number(authForm.Username, authForm.CountryCode); !ok {
					c.Response(false, fmt.Sprintf("verification:Phone number is invalid in your region %s", authForm.CountryCode), TokenResp{})
					return
				}
			}

			// check result through Email or Phone
			checkResult := object.CheckSigninCode(user, checkDest, authForm.Code, c.GetAcceptLanguage())
			if len(checkResult) != 0 {
				c.Response(false, checkResult, TokenResp{})
				return
			}

			// disable the verification code
			err := object.DisableVerificationCode(checkDest)
			if err != nil {
				c.Response(false, err.Error(), TokenResp{})
				return
			}
		}

		application, err := object.GetApplication(fmt.Sprintf("fireboom_%s", authForm.Application))
		if err != nil {
			c.Response(false, err.Error(), TokenResp{})
			return
		}

		if application == nil {
			c.Response(false, fmt.Sprintf("auth:The application: %s does not exist", authForm.Application), TokenResp{})
			return
		}

		resp = c.HandleLoggedIn(application, user, &authForm)

		record := object.NewRecord(c.Ctx)
		record.Organization = application.Organization
		record.User = user.Name
		util.SafeGoroutine(func() { object.AddRecord(record) })
	} else {
		if c.GetSessionUsername() != "" {
			// user already signed in to Casdoor, so let the user click the avatar button to do the quick sign-in
			application, err := object.GetApplication(fmt.Sprintf("fireboom_%s", authForm.Application))
			if err != nil {
				c.Response(false, err.Error(), TokenResp{})
				return
			}

			if application == nil {
				c.Response(false, fmt.Sprintf("auth:The application: %s does not exist", authForm.Application), TokenResp{})
				return
			}

			user := c.getCurrentUser()
			resp = c.HandleLoggedIn(application, user, &authForm)

			record := object.NewRecord(c.Ctx)
			record.Organization = application.Organization
			record.User = user.Name
			util.SafeGoroutine(func() { object.AddRecord(record) })
		} else {
			c.Response(false, fmt.Sprintf("auth:Unknown authentication type (not password or provider), form = %s", util.StructToJson(authForm)), TokenResp{})
			return
		}
	}

	c.Data["json"] = resp
	c.ServeJSON()
}
