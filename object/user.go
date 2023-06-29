package object

import (
	"casdoor/util"
	"fmt"
	"github.com/xorm-io/core"
	"strings"
)

type User struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100) index" json:"createdTime"`
	UpdatedTime string `xorm:"varchar(100)" json:"updatedTime"`

	Id                string `xorm:"varchar(100) index" json:"id"`
	Type              string `xorm:"varchar(100)" json:"type"`
	Password          string `xorm:"varchar(100)" json:"password"`
	PasswordSalt      string `xorm:"varchar(100)" json:"passwordSalt"`
	PasswordType      string `xorm:"varchar(100)" json:"passwordType"`
	DisplayName       string `xorm:"varchar(100)" json:"displayName"`
	Email             string `xorm:"varchar(100) index" json:"email"`
	EmailVerified     bool   `json:"emailVerified"`
	Phone             string `xorm:"varchar(20) index" json:"phone"`
	CountryCode       string `xorm:"varchar(6)" json:"countryCode"`
	SignupApplication string `xorm:"varchar(100)" json:"signupApplication"`
}

type Userinfo struct {
	Sub         string   `json:"sub"`
	Aud         string   `json:"aud"`
	Name        string   `json:"preferred_username,omitempty"`
	DisplayName string   `json:"name,omitempty"`
	Email       string   `json:"email,omitempty"`
	Avatar      string   `json:"picture,omitempty"`
	Address     string   `json:"address,omitempty"`
	Phone       string   `json:"phone,omitempty"`
	Groups      []string `json:"groups,omitempty"`
}

func (user *User) GetId() string {
	return fmt.Sprintf("%s/%s", user.Owner, user.Name)
}

func GetUser(id string) (*User, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getUser(owner, name)
}

func (user *User) GetCountryCode(countryCode string) string {
	if countryCode != "" {
		return countryCode
	}

	if user != nil && user.CountryCode != "" {
		return user.CountryCode
	}
	return ""
}

func AddUser(user *User) (bool, error) {
	var err error
	if user.Id == "" {
		user.Id = util.GenerateId()
	}

	if user.Owner == "" || user.Name == "" {
		return false, nil
	}

	// 查询组织，目前内置组织builtIn-->user.owner
	organization, _ := GetOrganizationByUser(user)
	if organization == nil {
		return false, nil
	}

	affected, err := adapter.Engine.Insert(user)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func GetUserByFields(organization string, field string) (*User, error) {
	// check username
	user, err := GetUserByField(organization, "name", field)
	if err != nil || user != nil {
		return user, err
	}

	// check email
	if strings.Contains(field, "@") {
		user, err = GetUserByField(organization, "email", field)
		if user != nil || err != nil {
			return user, err
		}
	}

	// check phone
	user, err = GetUserByField(organization, "phone", field)
	if user != nil || err != nil {
		return user, err
	}

	return nil, nil
}

func GetUserByField(organizationName string, field string, value string) (*User, error) {
	if field == "" || value == "" {
		return nil, nil
	}

	user := User{Owner: organizationName}
	existed, err := adapter.Engine.Where(fmt.Sprintf("%s=?", strings.ToLower(field)), value).Get(&user)
	if err != nil {
		return nil, err
	}

	if existed {
		return &user, nil
	} else {
		return nil, nil
	}
}

func HasUserByField(organizationName string, field string, value string) bool {
	user, err := GetUserByField(organizationName, field, value)
	if err != nil {
		panic(err)
	}
	return user != nil
}

func GetUserInfo(user *User, scope string, aud string, host string) *Userinfo {

	resp := Userinfo{
		Sub: user.Id,
		Aud: aud,
	}
	if strings.Contains(scope, "profile") {
		resp.Name = user.Name
		resp.DisplayName = user.DisplayName
	}
	if strings.Contains(scope, "email") {
		resp.Email = user.Email
	}
	if strings.Contains(scope, "phone") {
		resp.Phone = user.Phone
	}
	return &resp
}

func DeleteUser(user *User) (bool, error) {
	// Forced offline the user first
	_, err := DeleteSession(util.GetSessionId(user.Owner, user.Name, CasdoorApplication))
	if err != nil {
		return false, err
	}

	affected, err := adapter.Engine.ID(core.PK{user.Owner, user.Name}).Delete(&User{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func UpdateUser(id string, user *User, columns []string) (bool, error) {
	var err error
	owner, name := util.GetOwnerAndNameFromIdNoCheck(id)
	oldUser, err := getUser(owner, name)
	if err != nil {
		return false, err
	}
	if oldUser == nil {
		return false, nil
	}

	if user.Password == "***" {
		user.Password = oldUser.Password
	}

	if len(columns) == 0 {
		columns = []string{
			"owner", "display_name", "country_code", "signup_application",
		}
	}
	columns = append(columns, "name", "email", "phone", "country_code")

	affected, err := updateUser(oldUser, user, columns)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func updateUser(oldUser, user *User, columns []string) (int64, error) {
	session := adapter.Engine.NewSession()
	defer session.Close()

	session.Begin()

	affected, err := session.ID(core.PK{oldUser.Owner, oldUser.Name}).Cols(columns...).Update(user)
	if err != nil {
		session.Rollback()
		return affected, err
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		return 0, err
	}

	return affected, nil
}

func GetUserByPhone(owner string, phone string) (*User, error) {
	if owner == "" || phone == "" {
		return nil, nil
	}

	user := User{Owner: owner, Phone: phone}
	existed, err := adapter.Engine.Get(&user)
	if err != nil {
		return nil, err
	}

	if existed {
		return &user, nil
	} else {
		return nil, nil
	}
}

func GetUserByUserId(userId string) (*UserTokenInfo, error) {
	if userId == "" {
		return nil, fmt.Errorf("userId is blank")
	}
	_, name := util.GetOwnerAndNameFromId(userId)
	return GetUserTokenInfo(name)
}
