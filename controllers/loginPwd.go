package controllers

import (
	"casdoor/form"
	"casdoor/object"
	"fmt"
)

const PWD = "password"

func init() {
	// 通过密码登录
	authActionMap[PWD] = &authAction{
		login: func(authForm form.AuthForm) (user *object.User, err error) {
			application, err := object.GetApplication(fmt.Sprintf("fireboom_%s", authForm.Application))
			if err != nil {
				return
			}

			if application == nil {
				return nil, fmt.Errorf("auth:The application: %s does not exist", authForm.Application)
			}

			user, msg := object.CheckUserPassword(authForm.Organization, authForm.Username, authForm.Password)
			if msg != "" {
				return nil, fmt.Errorf(msg)
			}
			return
		},
	}

}
