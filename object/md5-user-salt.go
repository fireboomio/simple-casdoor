package object

import (
	"crypto/md5"
	"encoding/hex"
)

type Md5UserSaltCredManager struct{}

func init() {
	CredActionMap["md5"] = &Md5UserSaltCredManager{}
}

func getMd5(data []byte) []byte {
	hash := md5.Sum(data)
	return hash[:]
}

func getMd5HexDigest(s string) string {
	b := getMd5([]byte(s))
	res := hex.EncodeToString(b)
	return res
}

func (cm *Md5UserSaltCredManager) GetHashedPassword(password string, userSalt string) string {
	res := getMd5HexDigest(password)
	if userSalt != "" {
		res = getMd5HexDigest(res + userSalt)
	}
	return res
}

func (cm *Md5UserSaltCredManager) IsPasswordCorrect(plainPwd string, hashedPwd string, userSalt string) bool {
	return hashedPwd == cm.GetHashedPassword(plainPwd, userSalt)
}
