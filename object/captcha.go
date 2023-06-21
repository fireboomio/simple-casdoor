package object

import (
	"bytes"
	"github.com/dchest/captcha"
)

func GetCaptcha() (string, []byte, error) {
	id := captcha.NewLen(5)

	var buffer bytes.Buffer

	err := captcha.WriteImage(&buffer, id, 200, 80)
	if err != nil {
		return "", nil, err
	}

	return id, buffer.Bytes(), nil
}

func VerifyCaptcha(id string, digits string) bool {
	res := captcha.VerifyString(id, digits)

	return res
}
