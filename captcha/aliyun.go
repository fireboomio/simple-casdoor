package captcha

import (
	"casdoor/util"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const AliyunCaptchaVerifyUrl = "http://afs.aliyuncs.com"

type captchaSuccessResponse struct {
	Code int    `json:"Code"`
	Msg  string `json:"Msg"`
}

type captchaFailResponse struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

type AliyunCaptchaProvider struct{}

func NewAliyunCaptchaProvider() *AliyunCaptchaProvider {
	captcha := &AliyunCaptchaProvider{}
	return captcha
}

func contentEscape(str string) string {
	str = strings.Replace(str, " ", "%20", -1)
	str = url.QueryEscape(str)
	return str
}

func (captcha *AliyunCaptchaProvider) VerifyCaptcha(token, clientSecret string) (bool, error) {
	pathData, err := url.ParseQuery(token)
	if err != nil {
		return false, err
	}

	pathData["Action"] = []string{"AuthenticateSig"}
	pathData["Format"] = []string{"json"}
	pathData["SignatureMethod"] = []string{"HMAC-SHA1"}
	pathData["SignatureNonce"] = []string{strconv.FormatInt(time.Now().UnixNano(), 10)}
	pathData["SignatureVersion"] = []string{"1.0"}
	pathData["Timestamp"] = []string{time.Now().UTC().Format("2006-01-02T15:04:05Z")}
	pathData["Version"] = []string{"2018-01-12"}

	var keys []string
	for k := range pathData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortQuery := ""
	for _, k := range keys {
		sortQuery += k + "=" + contentEscape(pathData[k][0]) + "&"
	}
	sortQuery = strings.TrimSuffix(sortQuery, "&")

	stringToSign := fmt.Sprintf("GET&%s&%s", url.QueryEscape("/"), url.QueryEscape(sortQuery))

	signature := util.GetHmacSha1(clientSecret+"&", stringToSign)

	resp, err := http.Get(fmt.Sprintf("%s?%s&Signature=%s", AliyunCaptchaVerifyUrl, sortQuery, url.QueryEscape(signature)))
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return handleCaptchaResponse(body)
}

func handleCaptchaResponse(body []byte) (bool, error) {
	captchaResp := &captchaSuccessResponse{}
	err := json.Unmarshal(body, captchaResp)
	if err != nil {
		captchaFailResp := &captchaFailResponse{}
		err = json.Unmarshal(body, captchaFailResp)
		if err != nil {
			return false, err
		}

		return false, errors.New(captchaFailResp.Message)
	}

	return true, nil
}
