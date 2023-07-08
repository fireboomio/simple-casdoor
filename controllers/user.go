package controllers

import (
	"casdoor/object"
	"encoding/json"
	"fmt"
	"strings"
)

// AddUser
// @Title AddUser
// @Tag User API
// @Description add user
// @Param   name    			query   string  true   "名称"
// @Param   displayName 		query   string  false   "昵称"
// @Param   password    		query   string  true   "密码"
// @Param   phone	    		query   string  true   "电话号码"
// @Param   countryCode			query   string  false  "国际区号（默认CN）" Enums(CN, US, JP) default(CN)
// @Success 200 {object} controllers.Response 成功
// @router /add-user [post]
func (c *ApiController) AddUser() {
	var user object.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	msg := object.CheckUsername(user.Name)
	if msg != "" {
		c.ResponseError(msg)
		return
	}

	user.Owner = "builtIn"
	user.PasswordType = "plain"
	user.SignupApplication = "fireboom_builtIn"
	c.Data["json"] = wrapActionResponse(object.AddUser(&user))
	c.ServeJSON()
}

func (c *ApiController) UpdateUser() {
	id := c.Input().Get("id")
	columnsStr := c.Input().Get("columns")

	var user object.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if id == "" {
		id = c.GetSessionUsername()
		if id == "" {
			c.ResponseError("general:Missing parameter")
			return
		}
	}
	oldUser, err := object.GetUser(id)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if oldUser == nil {
		c.ResponseError(fmt.Sprintf("general:The user: %s doesn't exist", id))
		return
	}

	if oldUser.Owner == "builtIn" && oldUser.Name == "fireboom" && (user.Owner != "builtIn" || user.Name != "fireboom") {
		c.ResponseError("auth:Unauthorized operation")
		return
	}

	if msg := object.CheckUpdateUser(oldUser, &user); msg != "" {
		c.ResponseError(msg)
		return
	}

	columns := []string{}
	if columnsStr != "" {
		columns = strings.Split(columnsStr, ",")
	}

	affected, err := object.UpdateUser(id, &user, columns)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(affected)
	c.ServeJSON()
}

func (c *ApiController) DeleteUser() {
	var user object.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if user.Owner == "builtIn" && user.Name == "fireboom" {
		c.ResponseError("auth:Unauthorized operation")
		return
	}

	c.Data["json"] = wrapActionResponse(object.DeleteUser(&user))
	c.ServeJSON()
}

// GetUserByToken
// @Title User API
// @Description get user by token
// @Success 200 {object} controllers.UserResponse 成功
// @Router  /get-user [get]
func (c *ApiController) GetUserByToken() {
	userId := c.Ctx.Input.CruSession.Get("username")
	// 通过userId查询用户
	user, err := object.GetUserByUserId(userId.(string))
	if err != nil {
		c.ResponseError(err.Error())
	}
	c.ResponseToken(true, "", *user)
}
