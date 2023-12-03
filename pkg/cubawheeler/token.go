package cubawheeler

import (
	"time"

	"github.com/go-oauth2/oauth2/v4"
)

var _ oauth2.TokenInfo = &Token{}

type Token struct {
	AccessToken           string                     `json:"access_token" bson:"access_token"`
	AccessTokenCreatedAt  time.Time                  `json:"access_token_created_at" bson:"access_token_created_at"`
	AccessTokenExpiresAt  time.Duration              `json:"expires_at"  bson:"access_token_expires_at"`
	ClientID              string                     `json:"client_id,omitempty" bson:"client_id,omitempty"`
	Code                  string                     `json:"code,omitempty" bson:"code,omitempty"`
	CodeCreateAt          time.Time                  `json:"code_create_at,omitempty" bson:"code_create_at,omitempty"`
	CodeExpiresIn         time.Duration              `json:"code_expires_in,omitempty" bson:"code_expires_in,omitempty"`
	CodeChallenge         string                     `json:"code_challenge,omitempty" bson:"code_challenge,omitempty"`
	CodeChallengeMethod   oauth2.CodeChallengeMethod `json:"code_challenge_method,omitempty" bson:"code_challenge_method,omitempty"`
	RedirectUrl           string                     `json:"redirect_url,omitempty" bson:"redirect_url,omitempty"`
	RefreshToken          string                     `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	RefreshTokenCreateAt  time.Time                  `json:"refresh_token_create_at,omitempty" bson:"refresh_token_create_at,omitempty"`
	RefreshTokenExpiresIn time.Duration              `json:"refresh_token_expires_in,omitempty" bson:"refresh_token_expires_in,omitempty"`
	Scope                 string                     `json:"scope,omitempty" bson:"scope,omitempty"`
	UserID                string                     `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

// GetAccess implements oauth2.TokenInfo.
func (t *Token) GetAccess() string {
	return t.AccessToken
}

// GetAccessCreateAt implements oauth2.TokenInfo.
func (t *Token) GetAccessCreateAt() time.Time {
	return t.AccessTokenCreatedAt
}

// GetAccessExpiresIn implements oauth2.TokenInfo.
func (t *Token) GetAccessExpiresIn() time.Duration {
	return t.AccessTokenExpiresAt
}

// GetClientID implements oauth2.TokenInfo.
func (t *Token) GetClientID() string {
	return t.ClientID
}

// GetCode implements oauth2.TokenInfo.
func (t *Token) GetCode() string {
	return t.Code
}

// GetCodeChallenge implements oauth2.TokenInfo.
func (t *Token) GetCodeChallenge() string {
	return t.CodeChallenge
}

// GetCodeChallengeMethod implements oauth2.TokenInfo.
func (t *Token) GetCodeChallengeMethod() oauth2.CodeChallengeMethod {
	return t.CodeChallengeMethod
}

// GetCodeCreateAt implements oauth2.TokenInfo.
func (t *Token) GetCodeCreateAt() time.Time {
	return t.CodeCreateAt
}

// GetCodeExpiresIn implements oauth2.TokenInfo.
func (t *Token) GetCodeExpiresIn() time.Duration {
	return t.CodeExpiresIn
}

// GetRedirectURI implements oauth2.TokenInfo.
func (t *Token) GetRedirectURI() string {
	return t.RedirectUrl
}

// GetRefresh implements oauth2.TokenInfo.
func (t *Token) GetRefresh() string {
	return t.RefreshToken
}

// GetRefreshCreateAt implements oauth2.TokenInfo.
func (t *Token) GetRefreshCreateAt() time.Time {
	return t.RefreshTokenCreateAt
}

// GetRefreshExpiresIn implements oauth2.TokenInfo.
func (t *Token) GetRefreshExpiresIn() time.Duration {
	return t.RefreshTokenExpiresIn
}

// GetScope implements oauth2.TokenInfo.
func (t *Token) GetScope() string {
	return t.Scope
}

// GetUserID implements oauth2.TokenInfo.
func (t *Token) GetUserID() string {
	return t.UserID
}

// New implements oauth2.TokenInfo.
func (t *Token) New() oauth2.TokenInfo {
	return &Token{}
}

// SetAccess implements oauth2.TokenInfo.
func (t *Token) SetAccess(token string) {
	t.AccessToken = token
}

// SetAccessCreateAt implements oauth2.TokenInfo.
func (t *Token) SetAccessCreateAt(created time.Time) {
	t.AccessTokenCreatedAt = created
}

// SetAccessExpiresIn implements oauth2.TokenInfo.
func (t *Token) SetAccessExpiresIn(expiry time.Duration) {
	t.AccessTokenExpiresAt = expiry
}

// SetClientID implements oauth2.TokenInfo.
func (t *Token) SetClientID(client string) {
	t.ClientID = client
}

// SetCode implements oauth2.TokenInfo.
func (t *Token) SetCode(code string) {
	t.Code = code
}

// SetCodeChallenge implements oauth2.TokenInfo.
func (t *Token) SetCodeChallenge(challenge string) {
	t.CodeChallenge = challenge
}

// SetCodeChallengeMethod implements oauth2.TokenInfo.
func (t *Token) SetCodeChallengeMethod(method oauth2.CodeChallengeMethod) {
	t.CodeChallengeMethod = method
}

// SetCodeCreateAt implements oauth2.TokenInfo.
func (t *Token) SetCodeCreateAt(createAt time.Time) {
	t.CodeCreateAt = createAt
}

// SetCodeExpiresIn implements oauth2.TokenInfo.
func (t *Token) SetCodeExpiresIn(expiry time.Duration) {
	t.CodeExpiresIn = expiry
}

// SetRedirectURI implements oauth2.TokenInfo.
func (t *Token) SetRedirectURI(redirectURL string) {
	t.RedirectUrl = redirectURL
}

// SetRefresh implements oauth2.TokenInfo.
func (t *Token) SetRefresh(refresh string) {
	t.RefreshToken = refresh
}

// SetRefreshCreateAt implements oauth2.TokenInfo.
func (t *Token) SetRefreshCreateAt(createAt time.Time) {
	t.RefreshTokenCreateAt = createAt
}

// SetRefreshExpiresIn implements oauth2.TokenInfo.
func (t *Token) SetRefreshExpiresIn(expiry time.Duration) {
	t.RefreshTokenExpiresIn = expiry
}

// SetScope implements oauth2.TokenInfo.
func (t *Token) SetScope(scope string) {
	t.Scope = scope
}

// SetUserID implements oauth2.TokenInfo.
func (t *Token) SetUserID(userID string) {
	t.UserID = userID
}
