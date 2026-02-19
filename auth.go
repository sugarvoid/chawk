package chawk

import "time"

type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresIn    int64     `json:"expires_in"`
	Expiry       time.Time `json:"expiry"`
}

func (t *Token) IsExpired() bool {
	if t == nil || t.AccessToken == "" {
		return true
	}
	// 2 minute early for safety might not need this
	// I had it in the python version, but can't remember why
	return time.Now().Add(2 * time.Minute).After(t.Expiry)
}
