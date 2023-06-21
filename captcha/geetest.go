package captcha

import (
	"casdoor/util"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const GEETESTCaptchaVerifyUrl = "http://gcaptcha4.geetest.com/validate"

type GEETESTCaptchaProvider struct{}

func NewGEETESTCaptchaProvider() *GEETESTCaptchaProvider {
	captcha := &GEETESTCaptchaProvider{}
	return captcha
}

func (captcha *GEETESTCaptchaProvider) VerifyCaptcha(token, clientSecret string) (bool, error) {
	pathData, err := url.ParseQuery(token)
	if err != nil {
		return false, err
	}

	signToken := util.GetHmacSha256(clientSecret, pathData["lot_number"][0])

	formData := make(url.Values)
	formData["lot_number"] = []string{pathData["lot_number"][0]}
	formData["captcha_output"] = []string{pathData["captcha_output"][0]}
	formData["pass_token"] = []string{pathData["pass_token"][0]}
	formData["gen_time"] = []string{pathData["gen_time"][0]}
	formData["sign_token"] = []string{signToken}
	captchaId := pathData["captcha_id"][0]

	cli := http.Client{Timeout: time.Second * 5}
	resp, err := cli.PostForm(fmt.Sprintf("%s?captcha_id=%s", GEETESTCaptchaVerifyUrl, captchaId), formData)
	if err != nil || resp.StatusCode != 200 {
		return false, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	type captchaResponse struct {
		Result string `json:"result"`
		Reason string `json:"reason"`
	}
	captchaResp := &captchaResponse{}
	err = json.Unmarshal(body, captchaResp)
	if err != nil {
		return false, err
	}

	if captchaResp.Result == "success" {
		return true, nil
	}

	return false, errors.New(captchaResp.Reason)
}
