package object

import (
	"casdoor/cred"
	"casdoor/form"
	"casdoor/util"
	"fmt"
	"regexp"
	"unicode"
)

var (
	reWhiteSpace     *regexp.Regexp
	reFieldWhiteList *regexp.Regexp
)

func init() {
	reWhiteSpace, _ = regexp.Compile(`\s`)
	reFieldWhiteList, _ = regexp.Compile(`^[A-Za-z0-9]+$`)
}

func CheckUserPassword(organization string, username string, password string, lang string) (*User, string) {
	user, err := GetUserByFields(organization, username)
	if err != nil {
		panic(err)
	}

	if user == nil {
		return nil, fmt.Sprintf("general:The user: %s doesn't exist", util.GetId(organization, username))
	}

	if msg := CheckPassword(user, password, lang); msg != "" {
		return nil, msg
	}

	return user, ""
}

func CheckPassword(user *User, password string, lang string) string {
	organization, err := GetOrganizationByUser(user)
	if err != nil {
		panic(err)
	}

	if organization == nil {
		return "check:Organization does not exist"
	}

	passwordType := user.PasswordType
	if passwordType == "" {
		passwordType = "plain"
	}

	credManager := cred.GetCredManager(passwordType)
	if credManager != nil {
		if credManager.IsPasswordCorrect(password, user.Password, user.PasswordSalt) {
			return ""
		}
		return fmt.Sprintf("check: error password: %s", password)
	} else {
		return fmt.Sprintf(lang, "check:unsupported password type: %s", passwordType)
	}
}

func CheckUserSignup(application *Application, organization *Organization, form *form.AuthForm, lang string) string {
	if organization == nil {
		return "check:Organization does not exist"
	}

	if len(form.Username) <= 1 {
		return "check:Username must have at least 2 characters"
	}
	if unicode.IsDigit(rune(form.Username[0])) {
		return "check:Username cannot start with a digit"
	}
	if util.IsEmailValid(form.Username) {
		return "check:Username cannot be an email address"
	}
	if reWhiteSpace.MatchString(form.Username) {
		return "check:Username cannot contain white spaces"
	}

	if msg := CheckUsername(form.Username, lang); msg != "" {
		return msg
	}

	if HasUserByField(organization.Name, "name", form.Username) {
		return "check:Username already exists"
	}
	if HasUserByField(organization.Name, "email", form.Email) {
		return "check:Email already exists"
	}
	if HasUserByField(organization.Name, "phone", form.Phone) {
		return "check:Phone already exists"
	}

	if len(form.Password) <= 5 {
		return "check:Password must have at least 6 characters"
	}

	if form.Phone == "" {
		return "check:Phone cannot be empty"
	} else {
		if HasUserByField(organization.Name, "phone", form.Phone) {
			return "check:Phone already exists"
		} else if !util.IsPhoneValid(form.Phone, form.CountryCode) {
			return "check:Phone number is invalid"
		}
	}

	return ""
}

func CheckUsername(username string, lang string) string {
	if username == "" {
		return "check:Empty username."
	} else if len(username) > 39 {
		return "check:Username is too long (maximum is 39 characters)."
	}

	exclude, _ := regexp.Compile("^[\u0021-\u007E]+$")
	if !exclude.MatchString(username) {
		return ""
	}

	re, _ := regexp.Compile("^[a-zA-Z0-9]+((?:-[a-zA-Z0-9]+)|(?:_[a-zA-Z0-9]+))*$")
	if !re.MatchString(username) {
		return "check:The username may only contain alphanumeric characters, underlines or hyphens, cannot have consecutive hyphens or underlines, and cannot begin or end with a hyphen or underline."
	}

	return ""
}

func CheckUpdateUser(oldUser, user *User, lang string) string {
	if user.DisplayName == "" {
		return "user:Display name cannot be empty"
	}

	if oldUser.Name != user.Name {
		if msg := CheckUsername(user.Name, lang); msg != "" {
			return msg
		}
		if HasUserByField(user.Owner, "name", user.Name) {
			return "check:Username already exists"
		}
	}
	if oldUser.Email != user.Email {
		if HasUserByField(user.Name, "email", user.Email) {
			return "check:Email already exists"
		}
	}
	if oldUser.Phone != user.Phone {
		if HasUserByField(user.Owner, "phone", user.Phone) {
			return "check:Phone already exists"
		}
	}

	return ""
}
