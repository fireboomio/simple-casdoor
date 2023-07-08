package routers

import (
	"casdoor/object"
	"casdoor/util"
	"fmt"

	"github.com/beego/beego/context"
)

func AutoSigninFilter(ctx *context.Context) {
	// HTTP Bearer token like "Authorization: Bearer 123"
	accessToken := util.GetMaxLenStr(ctx.Input.Query("accessToken"), ctx.Input.Query("access_token"), ParseBearerToken(ctx))

	if accessToken != "" {
		token, err := object.GetTokenByAccessToken(accessToken)
		if err != nil {
			responseError(ctx, err.Error())
			return
		}

		if token == nil {
			responseError(ctx, "Access token doesn't exist")
			return
		}

		if util.IsTokenExpired(token.CreatedTime, token.ExpiresIn) {
			responseError(ctx, "Access token has expired")
			return
		}

		userId := util.GetId(token.Organization, token.User)
		application, err := object.GetApplicationByUserId(fmt.Sprintf("fireboom/%s", token.Application))
		if err != nil {
			panic(err)
		}

		setSessionUser(ctx, userId)
		setSessionOidc(ctx, token.Scope, application.ClientId)
		return
	} else {
		setSessionUser(ctx, "")
	}

	// "/page?clientId=123&clientSecret=456"
	userId := getUsernameByClientIdSecret(ctx)
	if userId != "" {
		setSessionUser(ctx, userId)
		return
	}

	// "/page?username=builtIn/fireboom&password=123"
	userId = ctx.Input.Query("username")
	password := ctx.Input.Query("password")
	if userId != "" && password != "" {
		owner, name := util.GetOwnerAndNameFromId(userId)
		_, msg := object.CheckUserPassword(owner, name, password)
		if msg != "" {
			responseError(ctx, msg)
			return
		}

		setSessionUser(ctx, userId)
		return
	}
}
