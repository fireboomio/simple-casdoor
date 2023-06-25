package object

import (
	"casdoor/util"
	"os"
)

func InitDb() {
	existed := initBuiltInOrganization()
	if !existed {
		initBuiltInProvider()
		initBuiltInUser()
		initBuiltInApplication()
		initBuiltInCert()
	}
}

// 初始化builtIn组织
func initBuiltInOrganization() bool {
	organization, err := getOrganization("fireboom", "builtIn")
	if err != nil {
		panic(err)
	}

	if organization != nil {
		return true
	}

	organization = &Organization{
		Owner:       "fireboom",
		Name:        "builtIn",
		CreatedTime: util.GetCurrentTime(),
		Languages:   []string{"en", "zh", "es", "fr", "de", "id", "ja", "ko", "ru", "vi", "pt"},
	}
	_, err = AddOrganization(organization)
	if err != nil {
		panic(err)
	}

	return false
}

// 若用户不存在则初始化 fireboom 用户
func initBuiltInUser() {
	user, err := getUser("builtIn", "fireboom")
	if err != nil {
		panic(err)
	}
	if user != nil {
		return
	}

	user = &User{
		Owner:             "builtIn",
		Name:              "fireboom",
		CreatedTime:       util.GetCurrentTime(),
		Id:                util.GenerateId(),
		Type:              "normal-user",
		Password:          "123",
		DisplayName:       "Admin",
		Email:             "admin@example.com",
		Phone:             "12345678910",
		CountryCode:       "CN",
		SignupApplication: "fireboom_builtIn",
	}
	_, err = AddUser(user)
	if err != nil {
		panic(err)
	}
}

func getUser(owner string, name string) (*User, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	user := User{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&user)
	if err != nil {
		return nil, err
	}

	if existed {
		return &user, nil
	} else {
		return nil, nil
	}
}

// 初始化fireboom应用
func initBuiltInApplication() {
	application, err := getApplication("fireboom", "fireboom_builtIn")
	if err != nil {
		panic(err)
	}

	if application != nil {
		return
	}

	application = &Application{
		Owner:                "fireboom",
		Name:                 "fireboom_builtIn",
		CreatedTime:          util.GetCurrentTime(),
		Organization:         "builtIn",
		RefreshExpireInHours: 168,
		ExpireInHours:        100,
	}
	_, err = AddApplication(application)
	if err != nil {
		panic(err)
	}
}

func getApplication(owner string, name string) (*Application, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	application := Application{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&application)
	if err != nil {
		return nil, err
	}
	if existed {
		return &application, nil
	}
	return nil, nil
}

// 初始化阿里云SMS provider
func initBuiltInProvider() {
	provider, err := GetProvider(util.GetId("fireboom", "provider_sms"))
	if err != nil {
		panic(err)
	}

	if provider != nil {
		return
	}

	provider = &Provider{
		Owner:       "fireboom",
		Name:        "provider_sms",
		CreatedTime: util.GetCurrentTime(),
		Category:    "SMS",
		Type:        "Aliyun SMS",
	}
	_, err = AddProvider(provider)
	if err != nil {
		panic(err)
	}
}

func readTokenFromFile() (string, string) {
	pemPath := "./object/token_jwt_key.pem"
	keyPath := "./object/token_jwt_key.key"
	pem, err := os.ReadFile(pemPath)
	if err != nil {
		return "", ""
	}
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return "", ""
	}
	return string(pem), string(key)
}

// 初始化私钥
func initBuiltInCert() {
	tokenJwtCertificate, tokenJwtPrivateKey := readTokenFromFile()
	cert, err := getCert("fireboom", "cert-builtIn")
	if err != nil {
		panic(err)
	}

	if cert != nil {
		return
	}

	cert = &Cert{
		Owner:           "fireboom",
		Name:            "cert-builtIn",
		CreatedTime:     util.GetCurrentTime(),
		DisplayName:     "builtIn Cert",
		Scope:           "JWT",
		Type:            "x509",
		CryptoAlgorithm: "RS256",
		BitSize:         4096,
		ExpireInYears:   20,
		Certificate:     tokenJwtCertificate,
		PrivateKey:      tokenJwtPrivateKey,
	}
	_, err = AddCert(cert)
	if err != nil {
		panic(err)
	}
}
