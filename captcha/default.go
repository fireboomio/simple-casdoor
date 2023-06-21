package captcha

import "casdoor/object"

type DefaultCaptchaProvider struct{}

func NewDefaultCaptchaProvider() *DefaultCaptchaProvider {
	captcha := &DefaultCaptchaProvider{}
	return captcha
}

func (captcha *DefaultCaptchaProvider) VerifyCaptcha(token, clientSecret string) (bool, error) {
	return object.VerifyCaptcha(clientSecret, token), nil
}
