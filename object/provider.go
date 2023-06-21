package object

import (
	"casdoor/util"
	"fmt"
)

type Provider struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk unique" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	Type         string `xorm:"varchar(100)" json:"type"`
	Category     string `xorm:"varchar(100)" json:"category"`
	ClientId     string `xorm:"varchar(100)" json:"clientId"`
	ClientSecret string `xorm:"varchar(2000)" json:"clientSecret"`
	SignName     string `xorm:"varchar(100)" json:"signName"`
	TemplateCode string `xorm:"varchar(100)" json:"templateCode"`
}

func getProvider(owner string, name string) (*Provider, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	provider := Provider{Name: name}
	existed, err := adapter.Engine.Get(&provider)
	if err != nil {
		return &provider, err
	}

	if existed {
		return &provider, nil
	} else {
		return nil, nil
	}
}

func GetProvider(id string) (*Provider, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getProvider(owner, name)
}

func GetCaptchaProviderByOwnerName(applicationId, lang string) (*Provider, error) {
	owner, name := util.GetOwnerAndNameFromId(applicationId)
	provider := Provider{Owner: owner, Name: name, Category: "Captcha"}
	existed, err := adapter.Engine.Get(&provider)
	if err != nil {
		return nil, err
	}

	if !existed {
		return nil, fmt.Errorf("provider:the provider: %s does not exist", applicationId)
	}

	return &provider, nil
}

func AddProvider(provider *Provider) (bool, error) {
	if provider.Owner == "" {
		provider.Owner = "fireboom"
	}
	affected, err := adapter.Engine.Insert(provider)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func UpdateProvider(id string, provider *Provider) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	if p, err := getProvider(owner, name); err != nil {
		return false, err
	} else if p == nil {
		return false, nil
	}

	affected, err := adapter.Engine.Table("provider").Where("owner=? and name=?", owner, name).Update(provider)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}
