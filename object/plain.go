package object

type PlainCredManager struct{}

func init() {
	CredActionMap["plain"] = &PlainCredManager{}
}

func (cm *PlainCredManager) GetHashedPassword(password string, userSalt string) string {
	return password
}

func (cm *PlainCredManager) IsPasswordCorrect(plainPwd string, hashedPwd string, userSalt string) bool {
	return hashedPwd == plainPwd
}
