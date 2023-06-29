package controllers

import (
	"casdoor/object"
	"casdoor/util"
	"time"

	"github.com/beego/beego"
	"github.com/beego/beego/logs"
)

type ApiController struct {
	beego.Controller
}

type SessionData struct {
	ExpireTime int64
}

func (c *ApiController) getCurrentUser() *object.User {
	var user *object.User
	var err error
	userId := c.GetSessionUsername()
	if userId == "" {
		user = nil
	} else {
		user, err = object.GetUser(userId)
		if err != nil {
			panic(err)
		}
	}
	return user
}

// GetSessionUsername ...
func (c *ApiController) GetSessionUsername() string {
	// check if user session expired
	sessionData := c.GetSessionData()

	if sessionData != nil &&
		sessionData.ExpireTime != 0 &&
		sessionData.ExpireTime < time.Now().Unix() {
		c.ClearUserSession()
		return ""
	}

	user := c.GetSession("username")
	if user == nil {
		return ""
	}

	return user.(string)
}

func (c *ApiController) GetSessionApplication() *object.Application {
	clientId := c.GetSession("aud")
	if clientId == nil {
		return nil
	}
	application, err := object.GetApplicationByClientId(clientId.(string))
	if err != nil {
		panic(err)
	}

	return application
}

func (c *ApiController) ClearUserSession() {
	c.SetSessionUsername("")
	c.SetSessionData(nil)
}

func (c *ApiController) GetSessionOidc() (string, string) {
	sessionData := c.GetSessionData()
	if sessionData != nil &&
		sessionData.ExpireTime != 0 &&
		sessionData.ExpireTime < time.Now().Unix() {
		c.ClearUserSession()
		return "", ""
	}
	scopeValue := c.GetSession("scope")
	audValue := c.GetSession("aud")
	var scope, aud string
	var ok bool
	if scope, ok = scopeValue.(string); !ok {
		scope = ""
	}
	if aud, ok = audValue.(string); !ok {
		aud = ""
	}
	return scope, aud
}

// SetSessionUsername ...
func (c *ApiController) SetSessionUsername(user string) {
	c.SetSession("username", user)
}

// GetSessionData ...
func (c *ApiController) GetSessionData() *SessionData {
	session := c.GetSession("SessionData")
	if session == nil {
		return nil
	}

	sessionData := &SessionData{}
	err := util.JsonToStruct(session.(string), sessionData)
	if err != nil {
		logs.Error("GetSessionData failed, error: %s", err)
		return nil
	}

	return sessionData
}

// SetSessionData ...
func (c *ApiController) SetSessionData(s *SessionData) {
	if s == nil {
		c.DelSession("SessionData")
		return
	}

	c.SetSession("SessionData", util.StructToJson(s))
}

func (c *ApiController) setExpireForSession() {
	timestamp := time.Now().Unix()
	timestamp += 3600 * 24
	c.SetSessionData(&SessionData{
		ExpireTime: timestamp,
	})
}

func wrapActionResponse(affected bool, e ...error) *Response {
	if len(e) != 0 && e[0] != nil {
		return &Response{Status: "error", Msg: e[0].Error()}
	} else if affected {
		return &Response{Status: "ok", Msg: "", Data: "Affected"}
	} else {
		return &Response{Status: "ok", Msg: "", Data: "UnAffected"}
	}
}

func wrapErrorResponse(err error) *Response {
	if err == nil {
		return &Response{Status: "ok", Msg: ""}
	} else {
		return &Response{Status: "error", Msg: err.Error()}
	}
}

func wrapErrorUserResponse(err error) *UserResponse {
	if err == nil {
		return &UserResponse{true, ""}
	} else {
		return &UserResponse{
			Success: false,
			Data:    err.Error(),
		}
	}
}

func (c *ApiController) Finish() {
	c.Controller.Finish()
}
