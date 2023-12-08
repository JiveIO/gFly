package facebook

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"app/core/oauth"
)

// Session stores data during the auth process with Facebook.
type Session struct {
	AuthURL     string
	AccessToken string
	ExpiresAt   time.Time
}

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the Facebook provider.
func (s Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New(oauth.NoAuthUrlErrorMessage)
	}
	return s.AuthURL, nil
}

// Authorize the session with Facebook and return the access token to be stored for future use.
func (s *Session) Authorize(provider oauth.Provider, params oauth.Params) (string, error) {
	p := provider.(*Provider)
	token, err := p.config.Exchange(oauth.ContextForClient(p.Client()), params.Get("code"))
	if err != nil {
		return "", err
	}

	if !token.Valid() {
		return "", errors.New("invalid token received from provider")
	}

	s.AccessToken = token.AccessToken
	s.ExpiresAt = token.Expiry
	return token.AccessToken, err
}

// Marshal the session into a string
func (s Session) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s Session) String() string {
	return s.Marshal()
}

// UnmarshalSession will unmarshal a JSON string into a session.
func (p *Provider) UnmarshalSession(data string) (oauth.Session, error) {
	sess := &Session{}
	err := json.NewDecoder(strings.NewReader(data)).Decode(sess)
	return sess, err
}
