package form

type AuthForm struct {
	Type string `json:"type"`

	Organization string `json:"organization"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`

	Application string `json:"application"`
	ClientId    string `json:"clientId"`
	Provider    string `json:"provider"`
	Code        string `json:"code"`
	Method      string `json:"method"`

	EmailCode   string `json:"emailCode"`
	PhoneCode   string `json:"phoneCode"`
	CountryCode string `json:"countryCode"`

	AutoSignin bool `json:"autoSignin"`

	CaptchaType  string `json:"captchaType"`
	CaptchaToken string `json:"captchaToken"`
	ClientSecret string `json:"clientSecret"`

	Passcode string `json:"passcode"`
}
