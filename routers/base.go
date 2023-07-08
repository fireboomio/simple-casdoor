package routers

import (
	"casdoor/object"
	"casdoor/util"
	"fmt"
	"strings"

	"github.com/beego/beego/context"
)

type Response struct {
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Data2  interface{} `json:"data2"`
}

func responseError(ctx *context.Context, error string, data ...interface{}) {
	resp := Response{Status: "error", Msg: error}
	switch len(data) {
	case 2:
		resp.Data2 = data[1]
		fallthrough
	case 1:
		resp.Data = data[0]
	}

	err := ctx.Output.JSON(resp, true, false)
	if err != nil {
		panic(err)
	}
}

func denyRequest(ctx *context.Context) {
	responseError(ctx, "auth:Unauthorized operation")
}

func getUsernameByClientIdSecret(ctx *context.Context) string {
	clientId, clientSecret, ok := ctx.Request.BasicAuth()
	if !ok {
		clientId = ctx.Input.Query("clientId")
		clientSecret = ctx.Input.Query("clientSecret")
	}

	if clientId == "" || clientSecret == "" {
		return ""
	}

	application, err := object.GetApplicationByClientId(clientId)
	if err != nil {
		panic(err)
	}

	if application == nil || application.ClientSecret != clientSecret {
		return ""
	}

	return fmt.Sprintf("fireboom/%s", application.Name)
}

func getSessionUser(ctx *context.Context) string {
	user := ctx.Input.CruSession.Get("username")
	if user == nil {
		return ""
	}

	return user.(string)
}

func setSessionUser(ctx *context.Context, user string) {
	err := ctx.Input.CruSession.Set("username", user)
	if err != nil {
		panic(err)
	}
	ctx.Input.CruSession.SessionRelease(ctx.ResponseWriter)
}

func setSessionExpire(ctx *context.Context, ExpireTime int64) {
	SessionData := struct{ ExpireTime int64 }{ExpireTime: ExpireTime}
	err := ctx.Input.CruSession.Set("SessionData", util.StructToJson(SessionData))
	if err != nil {
		panic(err)
	}
	ctx.Input.CruSession.SessionRelease(ctx.ResponseWriter)
}

func setSessionOidc(ctx *context.Context, scope string, aud string) {
	err := ctx.Input.CruSession.Set("scope", scope)
	if err != nil {
		panic(err)
	}
	err = ctx.Input.CruSession.Set("aud", aud)
	if err != nil {
		panic(err)
	}
	ctx.Input.CruSession.SessionRelease(ctx.ResponseWriter)
}

func ParseBearerToken(ctx *context.Context) string {
	header := ctx.Request.Header.Get("Authorization")
	tokens := strings.Split(header, " ")
	if len(tokens) != 2 {
		return ""
	}

	prefix := tokens[0]
	if prefix != "Bearer" {
		return ""
	}

	return tokens[1]
}
