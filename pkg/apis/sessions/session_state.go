package sessions

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/clock"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/encryption"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/logger"
	"github.com/pierrec/lz4"
	"github.com/vmihailenco/msgpack/v4"
)

// SessionState is used to store information about the currently authenticated user session
//Webjet propriatory cookie are implemented in "partial" pkg\apis\sessions\session_state_Webjet.go file
type SessionState struct {
	CreatedAt *time.Time `msgpack:"ca,omitempty"`
	ExpiresOn *time.Time `msgpack:"eo,omitempty"`

	AccessToken  string `msgpack:"at,omitempty"`
	IDToken      string `msgpack:"it,omitempty"`
	RefreshToken string `msgpack:"rt,omitempty"`

	Nonce []byte `msgpack:"n,omitempty"`

	Email             string   `msgpack:"e,omitempty"`
	User              string   `msgpack:"u,omitempty"`
	Groups            []string `msgpack:"g,omitempty"`
	PreferredUsername string   `msgpack:"pu,omitempty"`

	// Internal helpers, not serialized
	Clock clock.Clock `msgpack:"-"`
	Lock  Lock        `msgpack:"-"`
}

func (s *SessionState) ObtainLock(ctx context.Context, expiration time.Duration) error {
	if s.Lock == nil {
		s.Lock = &NoOpLock{}
	}
	return s.Lock.Obtain(ctx, expiration)
}

func (s *SessionState) RefreshLock(ctx context.Context, expiration time.Duration) error {
	if s.Lock == nil {
		s.Lock = &NoOpLock{}
	}
	return s.Lock.Refresh(ctx, expiration)
}

func (s *SessionState) ReleaseLock(ctx context.Context) error {
	if s.Lock == nil {
		s.Lock = &NoOpLock{}
	}
	return s.Lock.Release(ctx)
}

func (s *SessionState) PeekLock(ctx context.Context) (bool, error) {
	if s.Lock == nil {
		s.Lock = &NoOpLock{}
	}
	return s.Lock.Peek(ctx)
}

// CreatedAtNow sets a SessionState's CreatedAt to now
func (s *SessionState) CreatedAtNow() {
	now := s.Clock.Now()
	s.CreatedAt = &now
}

// SetExpiresOn sets an expiration
func (s *SessionState) SetExpiresOn(exp time.Time) {
	s.ExpiresOn = &exp
}

// ExpiresIn sets an expiration a certain duration from CreatedAt.
// CreatedAt will be set to time.Now if it is unset.
func (s *SessionState) ExpiresIn(d time.Duration) {
	if s.CreatedAt == nil {
		s.CreatedAtNow()
	}
	exp := s.CreatedAt.Add(d)
	s.ExpiresOn = &exp
}

// IsExpired checks whether the session has expired
func (s *SessionState) IsExpired() bool {
	if s.ExpiresOn != nil && !s.ExpiresOn.IsZero() && s.ExpiresOn.Before(s.Clock.Now()) {
		return true
	}
	return false
}

// Age returns the age of a session
func (s *SessionState) Age() time.Duration {
	if s.CreatedAt != nil && !s.CreatedAt.IsZero() {
		return s.Clock.Now().Truncate(time.Second).Sub(*s.CreatedAt)
	}
	return 0
}

// String constructs a summary of the session state
func (s *SessionState) String() string {
	o := fmt.Sprintf("Session{email:%s user:%s PreferredUsername:%s", s.Email, s.User, s.PreferredUsername)
	if s.AccessToken != "" {
		o += " token:true"
	}
	if s.IDToken != "" {
		o += " id_token:true"
	}
	if s.CreatedAt != nil && !s.CreatedAt.IsZero() {
		o += fmt.Sprintf(" created:%s", s.CreatedAt)
	}
	if s.ExpiresOn != nil && !s.ExpiresOn.IsZero() {
		o += fmt.Sprintf(" expires:%s", s.ExpiresOn)
	}
	if s.RefreshToken != "" {
		o += " refresh_token:true"
	}
	if len(s.Groups) > 0 {
		o += fmt.Sprintf(" groups:%v", s.Groups)
	}
	return o + "}"
}

func (s *SessionState) GetClaim(claim string) []string {
	if s == nil {
		return []string{}
	}
	switch claim {
	case "access_token":
		return []string{s.AccessToken}
	case "id_token":
		return []string{s.IDToken}
	case "created_at":
		return []string{s.CreatedAt.String()}
	case "expires_on":
		return []string{s.ExpiresOn.String()}
	case "refresh_token":
		return []string{s.RefreshToken}
	case "email":
		return []string{s.Email}
	case "user":
		return []string{s.User}
	case "groups":
		groups := make([]string, len(s.Groups))
		copy(groups, s.Groups)
		return groups
	case "preferred_username":
		return []string{s.PreferredUsername}
	default:
		return []string{}
	}
}

// CheckNonce compares the Nonce against a potential hash of it
func (s *SessionState) CheckNonce(hashed string) bool {
	return encryption.CheckNonce(s.Nonce, hashed)
}

// EncodeSessionState returns an encrypted, lz4 compressed, MessagePack encoded session
func (s *SessionState) EncodeSessionState(c encryption.Cipher, compress bool) ([]byte, error) {
	//Use propriatory Webjet smaller cookie
	if true {
		return s.EncodeSessionStateWebjet(c)
	}
	// logger.Printf("TRACE: SessionState: %+v", s)
	packed, err := msgpack.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("error marshalling session state to msgpack: %w", err)
	}

	if !compress {
		return c.Encrypt(packed)
	}

	compressed, err := lz4Compress(packed)
	if err != nil {
		return nil, err
	}
	// logger.Printf("TRACE: SessionState: %+v", compressed)
	return c.Encrypt(compressed)
}

// DecodeSessionState decodes a LZ4 compressed MessagePack into a Session State
func DecodeSessionState(data []byte, c encryption.Cipher, compressed bool) (*SessionState, error) {
	//Use propriatory Webjet smaller cookie
	if true {
		return DecodeSessionStateWebjet(data, c)
	}
	decrypted, err := c.Decrypt(data)
	if err != nil {
		return nil, fmt.Errorf("error decrypting the session state: %w", err)
	}

	packed := decrypted
	if compressed {
		packed, err = lz4Decompress(decrypted)
		if err != nil {
			return nil, err
		}
	}

	var ss SessionState
	err = msgpack.Unmarshal(packed, &ss)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling data to session state: %w", err)
	}

	return &ss, nil
}

// lz4Compress compresses with LZ4
//
// The Compress:Decompress ratio is 1:Many. LZ4 gives fastest decompress speeds
// at the expense of greater compression compared to other compression
// algorithms.
func lz4Compress(payload []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	zw := lz4.NewWriter(nil)
	zw.Header = lz4.Header{
		BlockMaxSize:     65536,
		CompressionLevel: 0,
	}
	zw.Reset(buf)

	reader := bytes.NewReader(payload)
	_, err := io.Copy(zw, reader)
	if err != nil {
		return nil, fmt.Errorf("error copying lz4 stream to buffer: %w", err)
	}
	err = zw.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing lz4 writer: %w", err)
	}

	compressed, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading lz4 buffer: %w", err)
	}

	return compressed, nil
}

// lz4Decompress decompresses with LZ4
func lz4Decompress(compressed []byte) ([]byte, error) {
	reader := bytes.NewReader(compressed)
	buf := new(bytes.Buffer)
	zr := lz4.NewReader(nil)
	zr.Reset(reader)
	_, err := io.Copy(buf, zr)
	if err != nil {
		return nil, fmt.Errorf("error copying lz4 stream to buffer: %w", err)
	}

	payload, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading lz4 buffer: %w", err)
	}

	return payload, nil
}

//The standard big cookie implementation is implemented in "main" pkg\apis\sessions\session_state.go file
//Old Webjet implementation of reduced  cookies s.userOrEmail(), s.AccessToken, s.ExpiresOn.Unix(), s.RefreshToken, s.Groups)
func (s *SessionState) EncodeSessionStateWebjet(c encryption.Cipher) ([]byte, error) {
	if c == nil || s.AccessToken == "" {
		return packAndEncrypt(s.userOrEmail(), c)
		//return s.userOrEmail(), nil
	}
	return s.encryptedString(c)
}

func (s *SessionState) userOrEmail() string {
	u := s.User
	if s.Email != "" {
		u = s.Email
	}
	return u
}

func (s *SessionState) encryptedString(c encryption.Cipher) ([]byte, error) {
	//var err error
	if c == nil {
		panic("error. missing cipher")
	}
	logger.Printf("TRACE: encryptedString SessionState: %+v", s)
	//content := fmt.Sprintf("%s:%s:%d:%s:%s", s.userOrEmail(), s.AccessToken, s.ExpiresOn.Unix(), s.RefreshToken, s.Groups)
	csvGroups := strings.Join(s.Groups[:], ",")
	logger.Printf("TRACE: encryptedString csvGroups: %v", csvGroups)
	content := fmt.Sprintf("%s:%s:%d:%s:%s", s.userOrEmail(), "", s.ExpiresOn.Unix(), "", csvGroups)
	logger.Printf("TRACE: encryptedString content: %v", content)
	return packAndEncrypt(content, c)
}
func packAndEncrypt(content string, c encryption.Cipher) ([]byte, error) {
	contentBytes, err := base64.URLEncoding.DecodeString(content)
	if err != nil {
		return nil, fmt.Errorf("error DecodeString:  %w %v", err, content)
	}
	return c.Encrypt(contentBytes)
	//return c.encryptString(content)
}
func DecodeSessionStateWebjet(data []byte, c encryption.Cipher) (s *SessionState, err error) {
	if c == nil {
		panic("error. missing cipher")
	}
	//dataString:=base64.URLEncoding.EncodeToString(data)
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

// Encrypt a value for use in a cookie
//From oauth2_proxy_Old\cookie\cookies.go  (c *Cipher) Encrypt
/*
func (c *encryption.Cipher) encryptString(value string) (string, error) {
	ciphertext := make([]byte, aes.BlockSize+len(value))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("failed to create initialization vector %s", err)
	}

	stream := cipher.NewCFBEncrypter(c.Block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(value))
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
*/

// Decrypt a value from a cookie to it's original string
//from cookie\cookies.go func (c *Cipher) Decrypt(s string) (string, error) {
// func (c *encryption.Cipher) DecryptBytes(value []byte) (string, error) {
// 	// s string
// 	// encrypted, err := base64.StdEncoding.DecodeString(s)
// 	// if err != nil {
// 	// 	return "", fmt.Errorf("failed to decrypt cookie value %s", err)
// 	// }
// 	encrypted:=value;

// 	if len(encrypted) < aes.BlockSize {
// 		return "", fmt.Errorf("encrypted cookie value should be "+
// 			"at least %d bytes, but is only %d bytes",
// 			aes.BlockSize, len(encrypted))
// 	}

// 	iv := encrypted[:aes.BlockSize]
// 	encrypted = encrypted[aes.BlockSize:]
// 	stream := cipher.NewCFBDecrypter(c.Block, iv)
// 	stream.XORKeyStream(encrypted, encrypted)

// 	return string(encrypted), nil
// }
