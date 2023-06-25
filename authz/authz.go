package authz

import (
	"casdoor/object"
	"casdoor/util"
	"strings"
)

func IsAllowed(subOwner string, subName string, method string, urlPath string) bool {

	user, err := object.GetUser(util.GetId(subOwner, subName))
	if err != nil {
		panic(err)
	}

	if user != nil && (subOwner == "builtIn") {
		return true
	}
	return false
}

func isAllowedInDemoMode(subOwner string, subName string, method string, urlPath string, objOwner string, objName string) bool {
	if method == "POST" {
		if strings.HasPrefix(urlPath, "/api/login") || urlPath == "/api/logout" || urlPath == "/api/signup" || urlPath == "/api/send-verification-code" || urlPath == "/api/send-email" || urlPath == "/api/verify-captcha" {
			return true
		} else if urlPath == "/api/update-user" {
			// Allow ordinary users to update their own information
			if subOwner == objOwner && subName == objName && !(subOwner == "builtIn" && subName == "admin") {
				return true
			}
			return false
		} else {
			return false
		}
	}

	// If method equals GET
	return true
}
