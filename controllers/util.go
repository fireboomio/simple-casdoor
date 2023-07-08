package controllers

import (
	"casdoor/object"
	"fmt"
)

// ResponseJsonData ...
func (c *ApiController) ResponseJsonData(resp *Response, data ...interface{}) {
	switch len(data) {
	case 2:
		resp.Data2 = data[1]
		fallthrough
	case 1:
		resp.Data = data[0]
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

// ResponseOk ...
func (c *ApiController) ResponseOk(data ...interface{}) {
	resp := &Response{Status: "ok"}
	c.ResponseJsonData(resp, data...)
}

func (c *ApiController) Response(success bool, msg string, data TokenResp) {
	resp := &UserResponse{success, msg, data}
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *ApiController) ResponseToken(success bool, msg string, data object.UserTokenInfo) {
	resp := &UserTokenResponse{success, msg, data}
	c.Data["json"] = resp
	c.ServeJSON()
}

// ResponseError ...
func (c *ApiController) ResponseError(error string, data ...interface{}) {
	resp := &Response{Status: "error", Msg: error}
	c.ResponseJsonData(resp, data...)
}

func (c *ApiController) RequireSignedIn() (string, bool) {
	userId := c.GetSessionUsername()
	if userId == "" {
		c.ResponseError("general:Please login first", "Please login first")
		return "", false
	}
	return userId, true
}

func (c *ApiController) RequireSignedInUser() (*object.User, bool) {
	userId, ok := c.RequireSignedIn()
	if !ok {
		return nil, false
	}

	user, err := object.GetUser(userId)
	if err != nil {
		panic(err)
	}

	if user == nil {
		c.ClearUserSession()
		c.ResponseError(fmt.Sprintf("general:The user: %s doesn't exist", userId))
		return nil, false
	}
	return user, true
}
