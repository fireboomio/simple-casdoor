package object

import (
	"casdoor/util"
	"errors"
	"fmt"
	"github.com/xorm-io/core"
	"math/rand"
	"strings"
	"time"
)

const (
	VerifyTypePhone = "phone"
	VerifyTypeEmail = "email"
)

const (
	VerificationSuccess = iota
	wrongCodeError
	noRecordError
	timeoutError
)

type VerificationRecord struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	RemoteAddr string `xorm:"varchar(100)"`
	Type       string `xorm:"varchar(10)"`
	User       string `xorm:"varchar(100) notnull"`
	Provider   string `xorm:"varchar(100) notnull"`
	Receiver   string `xorm:"varchar(100) notnull"`
	Code       string `xorm:"varchar(10) notnull"`
	Time       int64  `xorm:"notnull"`
	IsUsed     bool
}

type VerifyResult struct {
	Code int
	Msg  string
}

func IsAllowSend(user *User, remoteAddr, recordType string) error {
	var record VerificationRecord
	record.RemoteAddr = remoteAddr
	record.Type = recordType
	if user != nil {
		record.User = user.GetId() // owner/name
	}
	has, err := adapter.Engine.Desc("created_time").Get(&record)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	if has && now-record.Time < 60 {
		return errors.New("you can only send one code in 60s")
	}

	return nil
}

func CheckVerificationCode(dest, code string) *VerifyResult {
	record, err := getVerificationRecord(dest)
	if err != nil {
		panic(err)
	}

	if record == nil {
		return &VerifyResult{noRecordError, "verification:Code has not been sent yet!"}
	}

	var timeout int64 = 10
	if err != nil {
		panic(err)
	}

	now := time.Now().Unix()
	if now-record.Time > timeout*60 {
		return &VerifyResult{timeoutError, fmt.Sprintf("verification:You should verify your code in %d min!", timeout)}
	}

	if record.Code != code {
		return &VerifyResult{wrongCodeError, "verification:Wrong verification code!"}
	}

	return &VerifyResult{VerificationSuccess, ""}
}

func CheckSigninCode(user *User, dest, code string) string {
	result := CheckVerificationCode(dest, code)
	switch result.Code {
	case VerificationSuccess:
		return ""
	case wrongCodeError:
		return fmt.Sprintf("signIn code is wrong!")
	default:
		return result.Msg
	}
}

func SendVerificationCodeToPhone(user *User, provider *Provider, remoteAddr string, dest string) error {
	if provider == nil {
		return errors.New("please set a SMS provider first")
	}

	if err := IsAllowSend(user, remoteAddr, provider.Category); err != nil {
		return err
	}

	code := getRandomCode(6)
	if err := SendSms(provider, code, dest); err != nil {
		return err
	}

	if err := AddToVerificationRecord(user, provider, remoteAddr, provider.Category, dest, code); err != nil {
		return err
	}

	return nil
}

func AddToVerificationRecord(user *User, provider *Provider, remoteAddr, recordType, dest, code string) error {
	var record VerificationRecord
	record.RemoteAddr = remoteAddr
	record.Type = recordType
	if user != nil {
		record.User = user.GetId()
	}
	record.Owner = provider.Owner
	record.Name = util.GenerateId()
	record.CreatedTime = util.GetCurrentTime()

	record.Provider = provider.Name
	record.Receiver = dest
	record.Code = code
	record.Time = time.Now().Unix()
	record.IsUsed = false

	_, err := adapter.Engine.Insert(record)
	if err != nil {
		return err
	}

	return nil
}

func DisableVerificationCode(dest string) (err error) {
	record, err := getVerificationRecord(dest)
	if record == nil || err != nil {
		return
	}

	record.IsUsed = true
	_, err = adapter.Engine.ID(core.PK{record.Owner, record.Name}).AllCols().Update(record)
	return
}

func getVerificationRecord(dest string) (*VerificationRecord, error) {
	var record VerificationRecord
	record.Receiver = dest
	has, err := adapter.Engine.Desc("time").Where("is_used = false").Get(&record)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &record, nil
}

func GetVerifyType(username string) (verificationCodeType string) {
	if strings.Contains(username, "@") {
		return VerifyTypeEmail
	} else {
		return VerifyTypePhone
	}
}

var stdNums = []byte("0123456789")

func getRandomCode(length int) string {
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, stdNums[r.Intn(len(stdNums))])
	}
	return string(result)
}
