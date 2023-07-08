package object

import (
	"casdoor/util"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/xorm-io/core"
)

const (
	hourSeconds          = int(time.Hour / time.Second)
	InvalidRequest       = "invalid_request"
	InvalidClient        = "invalid_client"
	InvalidGrant         = "invalid_grant"
	UnauthorizedClient   = "unauthorized_client"
	UnsupportedGrantType = "unsupported_grant_type"
	InvalidScope         = "invalid_scope"
	EndpointError        = "endpoint_error"
)

type Code struct {
	Message string `xorm:"varchar(100)" json:"message"`
	Code    string `xorm:"varchar(100)" json:"code"`
}

type Token struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	Application  string `xorm:"varchar(100)" json:"application"`
	Organization string `xorm:"varchar(100)" json:"organization"`
	User         string `xorm:"varchar(100)" json:"user"`

	Code          string `xorm:"varchar(100) index" json:"code"`
	AccessToken   string `xorm:"mediumtext" json:"accessToken"`
	RefreshToken  string `xorm:"mediumtext" json:"refreshToken"`
	ExpiresIn     int    `json:"expiresIn"`
	Scope         string `xorm:"varchar(100)" json:"scope"`
	TokenType     string `xorm:"varchar(100)" json:"tokenType"`
	CodeChallenge string `xorm:"varchar(100)" json:"codeChallenge"`
	CodeIsUsed    bool   `json:"codeIsUsed"`
	CodeExpireIn  int64  `json:"codeExpireIn"`
}

type TokenWrapper struct {
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

type TokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type IntrospectionResponse struct {
	Active    bool     `json:"active"`
	Scope     string   `json:"scope,omitempty"`
	ClientId  string   `json:"client_id,omitempty"`
	Username  string   `json:"username,omitempty"`
	TokenType string   `json:"token_type,omitempty"`
	Exp       int64    `json:"exp,omitempty"`
	Iat       int64    `json:"iat,omitempty"`
	Nbf       int64    `json:"nbf,omitempty"`
	Sub       string   `json:"sub,omitempty"`
	Aud       []string `json:"aud,omitempty"`
	Iss       string   `json:"iss,omitempty"`
	Jti       string   `json:"jti,omitempty"`
}

func getToken(owner string, name string) (*Token, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	token := Token{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&token)
	if err != nil {
		return nil, err
	}

	if existed {
		return &token, nil
	}

	return nil, nil
}

func getTokenByCode(code string) (*Token, error) {
	token := Token{Code: code}
	existed, err := adapter.Engine.Get(&token)
	if err != nil {
		return nil, err
	}

	if existed {
		return &token, nil
	}

	return nil, nil
}

func updateUsedByCode(token *Token) bool {
	affected, err := adapter.Engine.Where("code=?", token.Code).Cols("code_is_used").Update(token)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func GetToken(id string) (*Token, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getToken(owner, name)
}

func (token *Token) GetId() string {
	return fmt.Sprintf("%s/%s", token.Owner, token.Name)
}

func UpdateToken(id string, token *Token) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	if t, err := getToken(owner, name); err != nil {
		return false, err
	} else if t == nil {
		return false, nil
	}

	affected, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(token)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func AddToken(token *Token) (bool, error) {
	affected, err := adapter.Engine.Insert(token)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func DeleteToken(token *Token) (bool, error) {
	affected, err := adapter.Engine.ID(core.PK{token.Owner, token.Name}).Delete(&Token{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func ExpireTokenByAccessToken(accessToken string) (bool, *Application, *Token, error) {
	token := Token{AccessToken: accessToken}
	existed, err := adapter.Engine.Get(&token)
	if err != nil {
		return false, nil, nil, err
	}

	if !existed {
		return false, nil, nil, nil
	}

	token.ExpiresIn = 0
	affected, err := adapter.Engine.ID(core.PK{token.Owner, token.Name}).Cols("expires_in").Update(&token)
	if err != nil {
		return false, nil, nil, err
	}

	application, err := getApplication(token.Owner, token.Application)
	if err != nil {
		return false, nil, nil, err
	}

	return affected != 0, application, &token, nil
}

func GetTokenByAccessToken(accessToken string) (*Token, error) {
	// Check if the accessToken is in the database
	token := Token{AccessToken: accessToken}
	existed, err := adapter.Engine.Get(&token)
	if err != nil {
		return nil, err
	}

	if !existed {
		return nil, nil
	}

	return &token, nil
}

func GetTokenByTokenAndApplication(token string, application string) (*Token, error) {
	tokenResult := Token{}
	existed, err := adapter.Engine.Where("(refresh_token = ? or access_token = ? ) and application = ?", token, token, application).Get(&tokenResult)
	if err != nil {
		return nil, err
	}

	if !existed {
		return nil, nil
	}

	return &tokenResult, nil
}

// PkceChallenge: base64-URL-encoded SHA256 hash of verifier, per rfc 7636
func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	challenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(sum[:])
	return challenge
}

// IsGrantTypeValid
// Check if grantType is allowed in the current application
// authorization_code is allowed by default
func IsGrantTypeValid(method string, grantTypes []string) bool {
	if method == "authorization_code" {
		return true
	}
	for _, m := range grantTypes {
		if m == method {
			return true
		}
	}
	return false
}

// GetAuthorizationCodeToken
// Authorization code flow
func GetAuthorizationCodeToken(application *Application, clientSecret string, code string, verifier string) (*Token, *TokenError, error) {
	if code == "" {
		return nil, &TokenError{
			Error:            InvalidRequest,
			ErrorDescription: "authorization code should not be empty",
		}, nil
	}

	token, err := getTokenByCode(code)
	if err != nil {
		return nil, nil, err
	}

	if token == nil {
		return nil, &TokenError{
			Error:            InvalidGrant,
			ErrorDescription: "authorization code is invalid",
		}, nil
	}
	if token.CodeIsUsed {
		// anti replay attacks
		return nil, &TokenError{
			Error:            InvalidGrant,
			ErrorDescription: "authorization code has been used",
		}, nil
	}

	if token.CodeChallenge != "" && pkceChallenge(verifier) != token.CodeChallenge {
		return nil, &TokenError{
			Error:            InvalidGrant,
			ErrorDescription: "verifier is invalid",
		}, nil
	}

	if application.Name != token.Application {
		return nil, &TokenError{
			Error:            InvalidGrant,
			ErrorDescription: "the token is for wrong application (client_id)",
		}, nil
	}

	if time.Now().Unix() > token.CodeExpireIn {
		// code must be used within 5 minutes
		return nil, &TokenError{
			Error:            InvalidGrant,
			ErrorDescription: "authorization code has expired",
		}, nil
	}
	return token, nil, nil
}

// GetTokenByUser
// Implicit flow
func GetTokenByUser(application *Application, user *User, scope string, host string) (*Token, error) {

	accessToken, refreshToken, tokenName, err := GenerateJwtToken(application, user, "", scope, host)
	if err != nil {
		return nil, err
	}

	token := &Token{
		Owner:        application.Owner,
		Name:         tokenName,
		CreatedTime:  util.GetCurrentTime(),
		Application:  application.Name,
		Organization: user.Owner,
		User:         user.Name,
		Code:         util.GenerateClientId(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    application.ExpireInHours * hourSeconds,
		Scope:        scope,
		TokenType:    "Bearer",
		CodeIsUsed:   true,
	}
	_, err = AddToken(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

type UserTokenInfo struct {
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}

func GetUserTokenInfo(name string) (*UserTokenInfo, error) {
	if name == "" {
		return nil, fmt.Errorf("name is null")
	}
	token := Token{
		Owner:       "fireboom",
		Application: "fireboom_builtIn",
		User:        name,
	}
	existed, err := adapter.Engine.Table("token").Desc("created_time").Get(&token)
	if err != nil {
		return nil, err
	}

	if existed {
		userInfo := UserTokenInfo{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresIn:    string(token.ExpiresIn),
			Username:     name,
		}
		return &userInfo, nil
	} else {
		return nil, nil
	}
}
