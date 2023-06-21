package controllers

import (
	"casdoor/object"
	"encoding/json"
)

func (c *ApiController) GetProvider() {
	id := c.Input().Get("id")

	provider, err := object.GetProvider(id)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(provider)
}

// UpdateProvider
// @Title UpdateProvider
// @Tag Provider API
// @Description update provider
// @Param   clientId    	query   string  true     "clientId"
// @Param   clientSecret    query   string  true     "clientSecret"
// @Param   signName    	query   string  true     "签名"
// @Param   templateCode    query   string  true     "模板代码"
// @Success 200 {object} controllers.Response "成功"
// @router /update-provider [post]
func (c *ApiController) UpdateProvider() {
	id := "fireboom/provider_sms"
	var provider object.Provider
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &provider)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	// 只更新相应的字段
	c.Data["json"] = wrapActionResponse(object.UpdateProvider(id, &provider))
	c.ServeJSON()
}
