package middleware

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc"
	sessionsapi "github.com/oauth2-proxy/oauth2-proxy/v7/pkg/apis/sessions"
	//	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/logger"
)

// TokenToSessionFunc takes a raw ID Token and converts it into a SessionState.
type TokenToSessionFunc func(ctx context.Context, token string) (*sessionsapi.SessionState, error)

// VerifyFunc takes a raw bearer token and verifies it returning the converted
// oidc.IDToken representation of the token.
type VerifyFunc func(ctx context.Context, token string) (*oidc.IDToken, error)

// CreateTokenToSessionFunc provides a handler that is a default implementation
// for converting a JWT into a session.
func CreateTokenToSessionFunc(verify VerifyFunc) TokenToSessionFunc {
	return func(ctx context.Context, token string) (*sessionsapi.SessionState, error) {
		var claims struct {
			Subject           string `json:"sub"`
			Email             string `json:"email"`
			Verified          *bool  `json:"email_verified"`
			PreferredUsername string `json:"preferred_username"`
			//Extra explicit claims not required, as we in Webjet mapped  Email,User/Subject (i.e. customerReferenceId), PreferredUsername(receiving as  preferred_username)(I.e. FirstName)
			// CustomerReferenceId string `json:"customerReferenceId"`
			// FirstName           string `json:"firstName"`
			// SignInName          string `json:"signInName"`
		}

		idToken, err := verify(ctx, token)
		if err != nil {
			return nil, err
		}

		if err := idToken.Claims(&claims); err != nil {
			return nil, fmt.Errorf("failed to parse bearer token claims: %v", err)
		}
		//		logger.Printf("TRACE: claims %+v", claims)
		if claims.Email == "" {
			claims.Email = claims.Subject
		}

		if claims.Verified != nil && !*claims.Verified {
			return nil, fmt.Errorf("email in id_token (%s) isn't verified", claims.Email)
		}

		newSession := &sessionsapi.SessionState{
			Email:             claims.Email,
			User:              claims.Subject,
			PreferredUsername: claims.PreferredUsername,
			AccessToken:       token,
			IDToken:           token,
			RefreshToken:      "",
			ExpiresOn:         &idToken.Expiry,
			//Extra explicit claims not required, as we in Webjet mapped  Email,User(i.e. customerReferenceId), PreferredUsername(receiving as  preferred_username)(I.e. FirstName)
			//CustomerReferenceId: claims.CustomerReferenceId,
			//FirstName:           claims.FirstName,
			// SignInName:          claims.SignInName,
		}

		return newSession, nil
	}
}
