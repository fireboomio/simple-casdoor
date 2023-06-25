package routers

import (
	"casdoor/authz"
	"casdoor/util"
	"encoding/json"
	"github.com/beego/beego/context"
	"net/http"
)

type Object struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

var allowAnonUrl = [...]string{"/api/login", "/api/send-verification-code"}

func isAllowAnonUrl(url string) bool {
	for i := 0; i < len(allowAnonUrl); i++ {
		if allowAnonUrl[i] == url {
			return true
		}
	}
	return false
}

func getUsername(ctx *context.Context) (username string) {
	defer func() {
		if r := recover(); r != nil {
			username = getUsernameByClientIdSecret(ctx)
		}
	}()

	username = ctx.Input.Session("username").(string)

	if username == "" {
		username = getUsernameByClientIdSecret(ctx)
	}

	return
}

func getSubject(ctx *context.Context) (string, string) {
	username := getUsername(ctx)
	if username == "" {
		return "anonymous", "anonymous"
	}
	return util.GetOwnerAndNameFromId(username)
}

func getObject(ctx *context.Context) (string, string) {
	method := ctx.Request.Method
	if method == http.MethodGet {
		// query == "?id=builtIn/fireboom"
		id := ctx.Input.Query("id")
		if id != "" {
			return util.GetOwnerAndNameFromId(id)
		}

		owner := ctx.Input.Query("owner")
		if owner != "" {
			return owner, ""
		}

		return "", ""
	} else {
		body := ctx.Input.RequestBody

		if len(body) == 0 {
			return ctx.Request.Form.Get("owner"), ctx.Request.Form.Get("name")
		}

		var obj Object
		err := json.Unmarshal(body, &obj)
		if err != nil {
			// panic(err)
			return "", ""
		}

		return obj.Owner, obj.Name
	}
}

func AuthzFilter(ctx *context.Context) {
	subOwner, subName := getSubject(ctx)
	method := ctx.Request.Method
	urlPath := ctx.Request.URL.Path

	var isAllowed bool
	if isAllowAnonUrl(urlPath) {
		isAllowed = true
	} else {
		isAllowed = authz.IsAllowed(subOwner, subName, method, urlPath)
	}
	if !isAllowed {
		denyRequest(ctx)
	}
}
