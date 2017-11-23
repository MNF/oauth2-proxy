package providers

import (
	"log"
	"oauth2_proxy/cookie"
	"strings"
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

const secret = "0123456789abcdefghijklmnopqrstuv"
const altSecret = "0000000000abcdefghijklmnopqrstuv"

func TestSessionStateSerialization(t *testing.T) {
	c, err := cookie.NewCipher([]byte(secret))
	assert.Equal(t, nil, err)
	c2, err := cookie.NewCipher([]byte(altSecret))
	assert.Equal(t, nil, err)
	s := &SessionState{
		Email:        "user@domain.com",
		AccessToken:  "token1234",
		ExpiresOn:    time.Now().Add(time.Duration(1) * time.Hour),
		RefreshToken: "refresh4321",
		Groups:       "test-group-1|test-group-2",
	}
	encoded, err := s.EncodeSessionState(c)
	assert.Equal(t, nil, err)
	log.Print(encoded)
	assert.Equal(t, 1, strings.Count(encoded, "=="))
	ss, err := DecodeSessionState(encoded, c)
	t.Logf("%#v", ss)
	assert.Equal(t, nil, err)
	assert.Equal(t, s.Email, ss.Email)
	assert.Equal(t, s.ExpiresOn.Unix(), ss.ExpiresOn.Unix())

	// ensure a different cipher can't decode properly (ie: it gets gibberish)
	ss, err = DecodeSessionState(encoded, c2)
	t.Logf("%#v", ss)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, s.Email, ss.Email)
	assert.NotEqual(t, s.ExpiresOn.Unix(), ss.ExpiresOn.Unix())
}

func TestSessionStateUserOrEmail(t *testing.T) {

	s := &SessionState{
		Email: "user@domain.com",
		User:  "just-user",
	}
	assert.Equal(t, "user@domain.com", s.userOrEmail())
	s.Email = ""
	assert.Equal(t, "just-user", s.userOrEmail())
}

func TestExpired(t *testing.T) {
	s := &SessionState{ExpiresOn: time.Now().Add(time.Duration(-1) * time.Minute)}
	assert.Equal(t, true, s.IsExpired())

	s = &SessionState{ExpiresOn: time.Now().Add(time.Duration(1) * time.Minute)}
	assert.Equal(t, false, s.IsExpired())

	s = &SessionState{}
	assert.Equal(t, false, s.IsExpired())
}
