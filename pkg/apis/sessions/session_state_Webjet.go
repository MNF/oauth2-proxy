package sessions

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/encryption"
	"github.com/vmihailenco/msgpack/v4"
)

//Old Webjet implementation of reduced  cookies s.userOrEmail(), s.AccessToken, s.ExpiresOn.Unix(), s.RefreshToken, s.Groups)
func (s *SessionState) EncodeSessionStateWebjet(c encryption.Cipher) ([]byte, error) {
	if c == nil || s.AccessToken == "" {
		return packAndEncrypt(s.userOrEmail(), c)
		//return s.userOrEmail(), nil
	}
	return s.EncryptedString(c)
}

func (s *SessionState) userOrEmail() string {
	u := s.User
	if s.Email != "" {
		u = s.Email
	}
	return u
}

func (s *SessionState) EncryptedString(c encryption.Cipher) ([]byte, error) {
	//var err error
	if c == nil {
		panic("error. missing cipher")
	}

	//content := fmt.Sprintf("%s:%s:%d:%s:%s", s.userOrEmail(), s.AccessToken, s.ExpiresOn.Unix(), s.RefreshToken, s.Groups)
	csvGroups := strings.Join(s.Groups[:], ",")
	content := fmt.Sprintf("%s:%s:%d:%s:%s", s.userOrEmail(), "", s.ExpiresOn.Unix(), "", csvGroups)
	return packAndEncrypt(content, c)
}
func packAndEncrypt(content string, c encryption.Cipher) ([]byte, error) {
	packed, err := msgpack.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf("error marshalling session state to msgpack: %w", err)
	}
	return c.Encrypt(packed)
	// content, err = c.Encrypt(packed)
	// if err != nil {
	// 	return "", err
	// }
	// return content, nil
}
func DecodeSessionStateWebjet(data []byte, c encryption.Cipher) (s *SessionState, err error) {
	if c == nil {
		panic("error. missing cipher")
	}
	vBytes, err := c.Decrypt(data)
	if err != nil {
		return nil, fmt.Errorf("error decrypting the session state: %w", err)
	}
	v := string(vBytes[:])
	chunks := strings.Split(v, ":")
	if len(chunks) == 1 {
		if strings.Contains(chunks[0], "@") {
			u := strings.Split(v, "@")[0]
			return &SessionState{Email: v, User: u}, nil
		}
		return &SessionState{User: v}, nil
	}

	if len(chunks) != 5 {
		err = fmt.Errorf("invalid number of fields (got %d expected 5)", len(chunks))
		return
	}

	s = &SessionState{}
	s.AccessToken = chunks[1]
	s.RefreshToken = chunks[3]
	if u := chunks[0]; strings.Contains(u, "@") {
		s.Email = u
		s.User = strings.Split(u, "@")[0]
	} else {
		s.User = u
	}
	if chunks[4] != "" {
		s.Groups = strings.Split(chunks[4], ",") //chunks[4]
	}
	ts, _ := strconv.Atoi(chunks[2])
	s.SetExpiresOn(time.Unix(int64(ts), 0)) // s.ExpiresOn = (time.Unix(int64(ts), 0)
	return
}
