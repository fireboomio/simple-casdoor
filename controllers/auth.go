package controllers

import (
	"casdoor/form"
	"casdoor/object"
	"casdoor/util"
	"encoding/json"
	"fmt"
	"sync"
)

type (
	authAction struct {
		login func(authForm form.AuthForm) (user *object.User, err error)
	}
	// TokenResp 返回 AccessToken 和 RefreshToken
	TokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)

var (
	lock          sync.RWMutex
	authActionMap map[string]*authAction
)

func init() {
	authActionMap = make(map[string]*authAction, 0)
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
// @Param password     	  query    string  true  "密码"
// @Param type     		  query    string  true  "类型：token"
// @Param loginType       query    string  true  "登录类型" Enums(sms, password)
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

		user, err = authActionMap[authForm.LoginType].login(authForm)
		if err != nil {
			c.Response(false, err.Error(), TokenResp{})
			return
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
