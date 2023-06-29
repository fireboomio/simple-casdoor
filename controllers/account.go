package controllers

import (
	"casdoor/form"
	"casdoor/object"
	"casdoor/util"
	"encoding/json"
	"fmt"
)

const (
	ResponseTypeLogin   = "login"
	ResponseTypeToken   = "token"
	ResponseTypeIdToken = "id_token"
)

type Response struct {
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Sub    string      `json:"sub"`
	Name   string      `json:"name"`
	Data   interface{} `json:"data"`
	Data2  interface{} `json:"data2"`
}

type UserResponse struct {
	Success bool      `json:"success"`
	Msg     string    `json:"msg"`
	Data    TokenResp `json:"data"`
}

type UserTokenResponse struct {
	Success bool                 `json:"success"`
	Msg     string               `json:"msg"`
	Data    object.UserTokenInfo `json:"data"`
}

type Captcha struct {
	Type          string `json:"type"`
	AppKey        string `json:"appKey"`
	Scene         string `json:"scene"`
	CaptchaId     string `json:"captchaId"`
	CaptchaImage  []byte `json:"captchaImage"`
	ClientId      string `json:"clientId"`
	ClientSecret  string `json:"clientSecret"`
	ClientId2     string `json:"clientId2"`
	ClientSecret2 string `json:"clientSecret2"`
	SubType       string `json:"subType"`
}

func (c *ApiController) Signup() {
	if c.GetSessionUsername() != "" {
		c.ResponseError("account:Please sign out first", c.GetSessionUsername())
		return
	}

	var authForm form.AuthForm
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &authForm)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	application, err := object.GetApplication(fmt.Sprintf("fireboom/%s", authForm.Application))
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	organization, err := object.GetOrganization(util.GetId("fireboom", authForm.Organization))
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	msg := object.CheckUserSignup(application, organization, &authForm, c.GetAcceptLanguage())
	if msg != "" {
		c.ResponseError(msg)
		return
	}

	var checkPhone string
	if authForm.Phone != "" {
		checkPhone, _ = util.GetE164Number(authForm.Phone, authForm.CountryCode)
		checkResult := object.CheckVerificationCode(checkPhone, authForm.PhoneCode, c.GetAcceptLanguage())
		if checkResult.Code != object.VerificationSuccess {
			c.ResponseError(checkResult.Msg)
			return
		}
	}

	id := util.GenerateId()

	username := authForm.Username

	user := &object.User{
		Owner:             authForm.Organization,
		Name:              username,
		CreatedTime:       util.GetCurrentTime(),
		Id:                id,
		Type:              "normal-user",
		Password:          authForm.Password,
		DisplayName:       authForm.Name,
		Email:             authForm.Email,
		Phone:             authForm.Phone,
		CountryCode:       authForm.CountryCode,
		SignupApplication: application.Name,
	}

	affected, err := object.AddUser(user)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if !affected {
		c.ResponseError("account:Failed to add user", util.StructToJson(user))
		return
	}

	err = object.DisableVerificationCode(authForm.Email)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	err = object.DisableVerificationCode(checkPhone)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	record := object.NewRecord(c.Ctx)
	record.Organization = application.Organization
	record.User = user.Name
	util.SafeGoroutine(func() { object.AddRecord(record) })

	userId := user.GetId()
	util.LogInfo(c.Ctx, "API: [%s] is signed up as new user", userId)

	c.ResponseOk(userId)
}

func (c *ApiController) Logout() {
	// https://openid.net/specs/openid-connect-rpinitiated-1_0-final.html
	accessToken := c.Input().Get("id_token_hint")
	redirectUri := c.Input().Get("post_logout_redirect_uri")

	user := c.GetSessionUsername()

	if accessToken == "" && redirectUri == "" {
		if user == "" {
			c.ResponseOk()
			return
		}

		c.ClearUserSession()
		owner, username := util.GetOwnerAndNameFromId(user)
		_, err := object.DeleteSessionId(util.GetSessionId(owner, username, object.CasdoorApplication), c.Ctx.Input.CruSession.SessionID())
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		util.LogInfo(c.Ctx, "API: [%s] logged out", user)

		application := c.GetSessionApplication()
		if application == nil || application.Name == "fireboom_builtIn" {
			c.ResponseOk(user)
			return
		}
		c.ResponseOk(user)
		return
	} else {
		if redirectUri == "" {
			c.ResponseError("general:Missing parameter: post_logout_redirect_uri")
			return
		}
		if accessToken == "" {
			c.ResponseError("general:Missing parameter: id_token_hint")
			return
		}

		affected, application, token, err := object.ExpireTokenByAccessToken(accessToken)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		if !affected {
			c.ResponseError("token:Token not found, invalid accessToken")
			return
		}

		if application == nil {
			c.ResponseError(fmt.Sprintf("auth:The application: %s does not exist"), token.Application)
			return
		}
	}
}

func (c *ApiController) GetAccount() {
	var err error
	user, ok := c.RequireSignedInUser()
	if !ok {
		return
	}

	organization, err := object.GetOrganizationByUser(user)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	resp := Response{
		Status: "ok",
		Sub:    user.Id,
		Name:   user.Name,
		Data:   user,
		Data2:  organization,
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *ApiController) GetUserinfo() {
	user, ok := c.RequireSignedInUser()
	if !ok {
		return
	}

	scope, aud := c.GetSessionOidc()
	host := c.Ctx.Request.Host
	userInfo := object.GetUserInfo(user, scope, aud, host)

	c.Data["json"] = userInfo
	c.ServeJSON()
}

func (c *ApiController) GetUserinfo2() {
	user, ok := c.RequireSignedInUser()
	if !ok {
		return
	}

	// this API is used by "Api URL" of Flarum's FoF Passport plugin
	// https://github.com/FriendsOfFlarum/passport
	type LaravelResponse struct {
		Id              string `json:"id"`
		Name            string `json:"name"`
		Email           string `json:"email"`
		EmailVerifiedAt string `json:"email_verified_at"`
		CreatedAt       string `json:"created_at"`
		UpdatedAt       string `json:"updated_at"`
	}

	response := LaravelResponse{
		Id:              user.Id,
		Name:            user.Name,
		Email:           user.Email,
		EmailVerifiedAt: user.CreatedTime,
		CreatedAt:       user.CreatedTime,
		UpdatedAt:       user.UpdatedTime,
	}

	c.Data["json"] = response
	c.ServeJSON()
}
