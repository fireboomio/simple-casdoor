package captcha

import "fmt"

type CaptchaProvider interface {
	VerifyCaptcha(token, clientSecret string) (bool, error)
}

func GetCaptchaProvider(captchaType string) CaptchaProvider {
	switch captchaType {
	case "Default":
		return NewDefaultCaptchaProvider()
	case "reCAPTCHA":
		return NewReCaptchaProvider()
	case "Aliyun Captcha":
		return NewAliyunCaptchaProvider()
	case "hCaptcha":
		return NewHCaptchaProvider()
	case "GEETEST":
		return NewGEETESTCaptchaProvider()
	case "Cloudflare Turnstile":
		return NewCloudflareTurnstileProvider()
	}

	return nil
}

func VerifyCaptchaByCaptchaType(captchaType, token, clientSecret string) (bool, error) {
	provider := GetCaptchaProvider(captchaType)
	if provider == nil {
		return false, fmt.Errorf("invalid captcha provider: %s", captchaType)
	}

	return provider.VerifyCaptcha(token, clientSecret)
}
