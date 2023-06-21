package captcha

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const CloudflareTurnstileVerifyUrl = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

type CloudflareTurnstileProvider struct{}

func NewCloudflareTurnstileProvider() *CloudflareTurnstileProvider {
	captcha := &CloudflareTurnstileProvider{}
	return captcha
}

func (captcha *CloudflareTurnstileProvider) VerifyCaptcha(token, clientSecret string) (bool, error) {
	reqData := url.Values{
		"secret":   {clientSecret},
		"response": {token},
	}
	resp, err := http.PostForm(CloudflareTurnstileVerifyUrl, reqData)
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
