package object

import (
	"casdoor/util"
	"fmt"
)

type Organization struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
}

func getOrganization(owner string, name string) (*Organization, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	organization := Organization{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&organization)
	if err != nil {
		return nil, err
	}

	if existed {
		return &organization, nil
	}
	return nil, nil
}

func GetOrganization(id string) (*Organization, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getOrganization(owner, name)
}

func UpdateOrganization(id string, organization *Organization) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	if org, err := getOrganization(owner, name); err != nil {
		return false, err
	} else if org == nil {
		return false, nil
	}

	if name == "builtIn" {
		organization.Name = name
	}

	if name != organization.Name {
		return false, fmt.Errorf("update organization failed: name %s not matched", name)
	}

	affected, err := adapter.Engine.Table("organization").Where("owner=? and name=?", owner, name).Update(organization)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func AddOrganization(organization *Organization) (bool, error) {
	affected, err := adapter.Engine.Insert(organization)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func GetOrganizationByUser(user *User) (*Organization, error) {
	return getOrganization("fireboom", user.Owner)
}
