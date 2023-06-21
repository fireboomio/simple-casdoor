package object

import (
	"strings"

	sender "github.com/casdoor/go-sms-sender"
)

func getSmsClient(provider *Provider) (sender.SmsClient, error) {
	var client sender.SmsClient
	var err error

	if provider.Type == sender.HuaweiCloud {
		client, err = sender.NewSmsClient(provider.Type, provider.ClientId, provider.ClientSecret, provider.SignName, provider.TemplateCode, "", "")
	} else {
		client, err = sender.NewSmsClient(provider.Type, provider.ClientId, provider.ClientSecret, provider.SignName, provider.TemplateCode, "")
	}
	if err != nil {
		return nil, err
	}

	return client, nil
}

func SendSms(provider *Provider, content string, phoneNumbers ...string) error {
	client, err := getSmsClient(provider)
	if err != nil {
		return err
	}

	if provider.Type == sender.Aliyun {
		for i, number := range phoneNumbers {
			phoneNumbers[i] = strings.TrimPrefix(number, "+86")
		}
	}

	params := map[string]string{}
	if provider.Type == sender.TencentCloud {
		params["0"] = content
	} else {
		params["code"] = content
	}

	err = client.SendMessage(params, phoneNumbers...)
	return err
}
