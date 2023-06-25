package form

import (
	"strings"
)

type VerificationForm struct {
	Dest          string `form:"dest"`
	Type          string `form:"type"`
	CountryCode   string `form:"countryCode"`
	ApplicationId string `form:"applicationId"`
	Method        string `form:"method"`
	CheckUser     string `form:"checkUser"`

	CaptchaType  string `form:"captchaType"`
	ClientSecret string `form:"clientSecret"`
	CaptchaToken string `form:"captchaToken"`
}

const (
	SendVerifyCode = 0
	VerifyCaptcha  = 1
)

func (form *VerificationForm) CheckParameter(checkType int, lang string) string {
	if checkType == SendVerifyCode {
		if form.Dest == "" {
			return "general:Missing parameter" + ": dest."
		}

		if !strings.Contains(form.ApplicationId, "_") {
			return "verification:Wrong parameter" + ": applicationId."
		}
	}

	if form.CaptchaType != "none" {
		if form.CaptchaToken == "" {
			return "general:Missing parameter" + ": captchaToken."
		}
		if form.ClientSecret == "" {
			return "general:Missing parameter" + ": clientSecret."
		}
	}

	return ""
}
