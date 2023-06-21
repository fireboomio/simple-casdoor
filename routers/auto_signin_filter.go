// Copyright 2021 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	// "/page?username=built-in/admin&password=123"
	userId = ctx.Input.Query("username")
	password := ctx.Input.Query("password")
	if userId != "" && password != "" {
		owner, name := util.GetOwnerAndNameFromId(userId)
		_, msg := object.CheckUserPassword(owner, name, password, "en")
		if msg != "" {
			responseError(ctx, msg)
			return
		}

		setSessionUser(ctx, userId)
		return
	}
}
