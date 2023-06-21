package object

import "casdoor/util"

type Application struct {
	Owner                string      `xorm:"varchar(100) notnull pk" json:"owner"`
	Name                 string      `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime          string      `xorm:"varchar(100)" json:"createdTime"`
	Organization         string      `xorm:"varchar(100)" json:"organization"`
	Providers            []*Provider `xorm:"mediumtext" json:"providers"`
	ExpireInHours        int         `json:"expireInHours"`
	RefreshExpireInHours int         `json:"refreshExpireInHours"`
	Cert                 string      `xorm:"varchar(100)" json:"cert"`
	ClientId             string      `xorm:"varchar(100)" json:"clientId"`
	ClientSecret         string      `xorm:"varchar(100)" json:"clientSecret"`
}

func AddApplication(application *Application) (bool, error) {
	if application.Owner == "" {
		application.Owner = "admin"
	}
	if application.Organization == "" {
		application.Organization = "built-in"
	}

	affected, err := adapter.Engine.Insert(application)
	if err != nil {
		return false, nil
	}

	return affected != 0, nil
}

func GetApplication(id string) (*Application, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getApplication(owner, name)
}

func GetApplicationByUserId(userId string) (application *Application, err error) {
	owner, name := util.GetOwnerAndNameFromId(userId)
	if owner == "fireboom" {
		application, err = getApplication("fireboom", name)
		return
	}

	user, err := GetUser(userId)
	if err != nil {
		return nil, err
	}
	application, err = GetApplicationByUser(user)
	return
}

func GetApplicationByUser(user *User) (*Application, error) {
	if user.SignupApplication != "" {
		return getApplication("fireboom", user.SignupApplication)
	} else {
		return GetApplicationByOrganizationName(user.Owner)
	}
}

func GetApplicationByOrganizationName(organization string) (*Application, error) {
	application := Application{}
	existed, err := adapter.Engine.Where("organization=?", organization).Get(&application)
	if err != nil {
		return nil, nil
	}

	if existed {
		return &application, nil
	} else {
		return nil, nil
	}
}

func GetApplicationByClientId(clientId string) (*Application, error) {
	application := Application{}
	existed, err := adapter.Engine.Where("client_id=?", clientId).Get(&application)
	if err != nil {
		return nil, err
	}

	if existed {
		return &application, nil
	} else {
		return nil, nil
	}
}
