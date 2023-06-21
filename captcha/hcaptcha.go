package captcha

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const HCaptchaVerifyUrl = "https://hcaptcha.com/siteverify"

type HCaptchaProvider struct{}

func NewHCaptchaProvider() *HCaptchaProvider {
	captcha := &HCaptchaProvider{}
	return captcha
}

func (captcha *HCaptchaProvider) VerifyCaptcha(token, clientSecret string) (bool, error) {
	reqData := url.Values{
		"secret":   {clientSecret},
		"response": {token},
	}
	resp, err := http.PostForm(HCaptchaVerifyUrl, reqData)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	type captchaResponse struct {
		Success    bool     `json:"success"`
		ErrorCodes []string `json:"error-codes"`
	}
	captchaResp := &captchaResponse{}
	err = json.Unmarshal(body, captchaResp)
	if err != nil {
		return false, err
	}

	if len(captchaResp.ErrorCodes) > 0 {
		return false, errors.New(strings.Join(captchaResp.ErrorCodes, ","))
	}

	return captchaResp.Success, nil
}
