package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/apis/sessions"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func testAzureProvider(hostname string) *AzureProvider {
	p := NewAzureProvider(
		&ProviderData{
			ProviderName:      "",
			LoginURL:          &url.URL{},
			RedeemURL:         &url.URL{},
			ProfileURL:        &url.URL{},
			ValidateURL:       &url.URL{},
			ProtectedResource: &url.URL{},
			Scope:             ""})

	if hostname != "" {
		updateURL(p.Data().LoginURL, hostname)
		updateURL(p.Data().RedeemURL, hostname)
		updateURL(p.Data().ProfileURL, hostname)
		updateURL(p.Data().ValidateURL, hostname)
		updateURL(p.Data().ProtectedResource, hostname)
	}
	return p
}

func TestNewAzureProvider(t *testing.T) {
	g := NewWithT(t)

	// Test that defaults are set when calling for a new provider with nothing set
	providerData := NewAzureProvider(&ProviderData{}).Data()
	g.Expect(providerData.ProviderName).To(Equal("Azure"))
	g.Expect(providerData.LoginURL.String()).To(Equal("https://login.microsoftonline.com/common/oauth2/authorize"))
	g.Expect(providerData.RedeemURL.String()).To(Equal("https://login.microsoftonline.com/common/oauth2/token"))
	g.Expect(providerData.ProfileURL.String()).To(Equal("https://graph.microsoft.com/v1.0/me"))
	g.Expect(providerData.ValidateURL.String()).To(Equal("https://graph.microsoft.com/v1.0/me"))
	g.Expect(providerData.Scope).To(Equal("openid"))
}

type mockTransport struct {
	params map[string]string
}

func TestAzureProviderOverrides(t *testing.T) {
	p := NewAzureProvider(
		&ProviderData{
			LoginURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/oauth/auth"},
			RedeemURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/oauth/token"},
			ProfileURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/oauth/profile"},
			ValidateURL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/oauth/tokeninfo"},
			ProtectedResource: &url.URL{
				Scheme: "https",
				Host:   "example.com"},
			Scope: "profile"})
	assert.NotEqual(t, nil, p)
	assert.Equal(t, "Azure", p.Data().ProviderName)
	assert.Equal(t, "https://example.com/oauth/auth",
		p.Data().LoginURL.String())
	assert.Equal(t, "https://example.com/oauth/token",
		p.Data().RedeemURL.String())
	assert.Equal(t, "https://example.com/oauth/profile",
		p.Data().ProfileURL.String())
	assert.Equal(t, "https://example.com/oauth/tokeninfo",
		p.Data().ValidateURL.String())
	assert.Equal(t, "https://example.com",
		p.Data().ProtectedResource.String())
	assert.Equal(t, "profile", p.Data().Scope)
}

func TestAzureSetTenant(t *testing.T) {
	p := testAzureProvider("")
	p.Configure("example")
	assert.Equal(t, "Azure", p.Data().ProviderName)
	assert.Equal(t, "example", p.Tenant)
	assert.Equal(t, "https://login.microsoftonline.com/example/oauth2/authorize",
		p.Data().LoginURL.String())
	assert.Equal(t, "https://login.microsoftonline.com/example/oauth2/token",
		p.Data().RedeemURL.String())
	assert.Equal(t, "https://graph.microsoft.com/v1.0/me",
		p.Data().ProfileURL.String())
	assert.Equal(t, "https://graph.microsoft.com",
		p.Data().ProtectedResource.String())
	assert.Equal(t, "https://graph.microsoft.com/v1.0/me", p.Data().ValidateURL.String())
	assert.Equal(t, "openid", p.Data().Scope)
}

func testAzureBackend(payload string) *httptest.Server {
	path := "/v1.0/me"

	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if (r.URL.Path != path) && r.Method != http.MethodPost {
				w.WriteHeader(404)
			} else if r.Method == http.MethodPost && r.Body != nil {
				w.WriteHeader(200)
				w.Write([]byte(payload))
			} else if !IsAuthorizedInHeader(r.Header) {
				w.WriteHeader(403)
			} else {
				w.WriteHeader(200)
				w.Write([]byte(payload))
			}
		}))
}

func TestAzureProviderGetEmailAddress(t *testing.T) {
	b := testAzureBackend(`{ "mail": "user@windows.net" }`)
	defer b.Close()

	bURL, _ := url.Parse(b.URL)
	p := testAzureProvider(bURL.Host)

	session := CreateAuthorizedSession()
	email, err := p.GetEmailAddress(context.Background(), session)
	assert.Equal(t, nil, err)
	assert.Equal(t, "user@windows.net", email)
}

func TestAzureProviderGetEmailAddressMailNull(t *testing.T) {
	b := testAzureBackend(`{ "mail": null, "otherMails": ["user@windows.net", "altuser@windows.net"] }`)
	defer b.Close()

	bURL, _ := url.Parse(b.URL)
	p := testAzureProvider(bURL.Host)

	session := CreateAuthorizedSession()
	email, err := p.GetEmailAddress(context.Background(), session)
	assert.Equal(t, nil, err)
	assert.Equal(t, "user@windows.net", email)
}

func TestAzureProviderGetEmailAddressGetUserPrincipalName(t *testing.T) {
	b := testAzureBackend(`{ "mail": null, "otherMails": [], "userPrincipalName": "user@windows.net" }`)
	defer b.Close()

	bURL, _ := url.Parse(b.URL)
	p := testAzureProvider(bURL.Host)

	session := CreateAuthorizedSession()
	email, err := p.GetEmailAddress(context.Background(), session)
	assert.Equal(t, nil, err)
	assert.Equal(t, "user@windows.net", email)
}

func TestAzureProviderGetEmailAddressFailToGetEmailAddress(t *testing.T) {
	b := testAzureBackend(`{ "mail": null, "otherMails": [], "userPrincipalName": null }`)
	defer b.Close()

	bURL, _ := url.Parse(b.URL)
	p := testAzureProvider(bURL.Host)

	session := CreateAuthorizedSession()
	email, err := p.GetEmailAddress(context.Background(), session)
	assert.Equal(t, "type assertion to string failed", err.Error())
	assert.Equal(t, "", email)
}

func TestAzureProviderGetEmailAddressEmptyUserPrincipalName(t *testing.T) {
	b := testAzureBackend(`{ "mail": null, "otherMails": [], "userPrincipalName": "" }`)
	defer b.Close()

	bURL, _ := url.Parse(b.URL)
	p := testAzureProvider(bURL.Host)

	session := CreateAuthorizedSession()
	email, err := p.GetEmailAddress(context.Background(), session)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", email)
}

func TestAzureProviderGetEmailAddressIncorrectOtherMails(t *testing.T) {
	b := testAzureBackend(`{ "mail": null, "otherMails": "", "userPrincipalName": null }`)
	defer b.Close()

	bURL, _ := url.Parse(b.URL)
	p := testAzureProvider(bURL.Host)

	session := CreateAuthorizedSession()
	email, err := p.GetEmailAddress(context.Background(), session)
	assert.Equal(t, "type assertion to string failed", err.Error())
	assert.Equal(t, "", email)
}

func TestAzureProviderRedeemReturnsIdToken(t *testing.T) {
	b := testAzureBackend(`{ "id_token": "testtoken1234", "expires_on": "1136239445", "refresh_token": "refresh1234" }`)
	defer b.Close()
	timestamp, err := time.Parse(time.RFC3339, "2006-01-02T22:04:05Z")
	assert.Equal(t, nil, err)

	bURL, _ := url.Parse(b.URL)
	p := testAzureProvider(bURL.Host)
	p.Data().RedeemURL.Path = "/common/oauth2/token"
	s, err := p.Redeem(context.Background(), "https://localhost", "1234")
	assert.Equal(t, nil, err)
	assert.Equal(t, "testtoken1234", s.IDToken)
	assert.Equal(t, timestamp, s.ExpiresOn.UTC())
	assert.Equal(t, "refresh1234", s.RefreshToken)
}

func TestAzureProviderProtectedResourceConfigured(t *testing.T) {
	p := testAzureProvider("")
	p.ProtectedResource, _ = url.Parse("http://my.resource.test")
	result := p.GetLoginURL("https://my.test.app/oauth", "")
	assert.Contains(t, result, "resource="+url.QueryEscape("http://my.resource.test"))
}

func TestAzureProviderGetsTokensInRedeem(t *testing.T) {
	b := testAzureBackend(`{ "access_token": "some_access_token", "refresh_token": "some_refresh_token", "expires_on": "1136239445", "id_token": "some_id_token" }`)
	defer b.Close()
	timestamp, _ := time.Parse(time.RFC3339, "2006-01-02T22:04:05Z")
	bURL, _ := url.Parse(b.URL)
	p := testAzureProvider(bURL.Host)

	session, err := p.Redeem(context.Background(), "http://redirect/", "code1234")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, session, nil)
	assert.Equal(t, "some_access_token", session.AccessToken)
	assert.Equal(t, "some_refresh_token", session.RefreshToken)
	assert.Equal(t, "some_id_token", session.IDToken)
	assert.Equal(t, timestamp, session.ExpiresOn.UTC())
}

func TestAzureProviderNotRefreshWhenNotExpired(t *testing.T) {
	p := testAzureProvider("")

	expires := time.Now().Add(time.Duration(1) * time.Hour)
	session := &sessions.SessionState{AccessToken: "some_access_token", RefreshToken: "some_refresh_token", IDToken: "some_id_token", ExpiresOn: &expires}
	refreshNeeded, err := p.RefreshSessionIfNeeded(context.Background(), session)
	assert.Equal(t, nil, err)
	assert.False(t, refreshNeeded)
}

func TestAzureProviderRefreshWhenExpired(t *testing.T) {
	b := testAzureBackend(`{ "access_token": "new_some_access_token", "refresh_token": "new_some_refresh_token", "expires_on": "32693148245", "id_token": "new_some_id_token" }`)
	defer b.Close()
	timestamp, _ := time.Parse(time.RFC3339, "3006-01-02T22:04:05Z")
	bURL, _ := url.Parse(b.URL)
	p := testAzureProvider(bURL.Host)

	expires := time.Now().Add(time.Duration(-1) * time.Hour)
	session := &sessions.SessionState{AccessToken: "some_access_token", RefreshToken: "some_refresh_token", IDToken: "some_id_token", ExpiresOn: &expires}
	_, err := p.RefreshSessionIfNeeded(context.Background(), session)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, session, nil)
	assert.Equal(t, "new_some_access_token", session.AccessToken)
	assert.Equal(t, "new_some_refresh_token", session.RefreshToken)
	assert.Equal(t, "new_some_id_token", session.IDToken)
	assert.Equal(t, timestamp, session.ExpiresOn.UTC())
}

/*TODO Convert TESTING GetGroups to testing addGroupsToSession
//Previous PermittedGroups are now allowed_groups
func newMockTransport(params map[string]string) http.RoundTripper {
	return &mockTransport{params}
}
func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("Starting Round Tripper")
	// Create mocked http.Response
	response := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
	}
	response.Header.Set("Content-Type", "application/json")

	//url := req.URL
	full_request := req.URL.Path
	if req.URL.RawQuery != "" {
		full_request += "?" + req.URL.RawQuery
	}
	var err error
	if value, ok := t.params[full_request]; ok {
		if req.Header.Get("Authorization") != "Bearer imaginary_access_token" {
			response.StatusCode = http.StatusForbidden
			err = fmt.Errorf("got 403. Bearer token '%v' is not correct", req.Header.Get("Authorization"))
		} else {
			response.StatusCode = http.StatusOK
			response.Body = ioutil.NopCloser(strings.NewReader(value))
			err = nil
		}

	} else {
		response.StatusCode = http.StatusNotFound
		err = fmt.Errorf("got 404. Requested path '%v' is not found", full_request)
	}

	return response, err
}
func TestAzureProviderNoGroups(t *testing.T) {
	params := map[string]string{}
		//path_group: payload_group_empty}

	http.DefaultClient.Transport = newMockTransport(params)

	p := testAzureProvider("")

	session := &sessions.SessionState{
		AccessToken: "imaginary_access_token",
		IDToken:     "imaginary_IDToken_token"}

	p.addGroupsToSession(session)//
	http.DefaultClient.Transport = nil
	assert.Equal(t, 0, len(session.Groups))
}

func TestAzureProviderWrongRequestGroups(t *testing.T) {
	params := map[string]string{
		path_group_wrong: payload_group_part_1}
	http.DefaultClient.Transport = newMockTransport(params)
	log.Printf("Def %#v\n\n", http.DefaultClient.Transport)

	p := testAzureProvider("")

	session := &sessions.SessionState{
		AccessToken: "imaginary_access_token",
		IDToken:     "imaginary_IDToken_token"}

	groups, _ := p.GetGroups(session, "")
	http.DefaultClient.Transport = nil

	//assert.NotEqual(t, nil, err)
	assert.Equal(t, "", groups)
}

func TestAzureProviderMultiRequestsGroups(t *testing.T) {
	params := map[string]string{
		path_group:      payload_group_part_1,
		path_group_next: payload_group_part_2}
	http.DefaultClient.Transport = newMockTransport(params)

	p := testAzureProvider("")

	session := &sessions.SessionState{
		AccessToken: "imaginary_access_token",
		IDToken:     "imaginary_IDToken_token"}

	groups, err := p.GetGroups(session, "")
	http.DefaultClient.Transport = nil

	assert.Equal(t, nil, err)
	assert.Equal(t, "test-group-1|test-group-2|test-group-3", groups)
}

func TestAzureEmptyPermittedGroups(t *testing.T) {
	p := testAzureProvider("")

	session := &sessions.SessionState{
		AccessToken: "imaginary_access_token",
		IDToken:     "imaginary_IDToken_token",
		Groups:      "no one|cares|non existing|groups"}
	result := p.ValidateGroup(session)

	assert.Equal(t, true, result)
}

func TestAzureWrongPermittedGroups(t *testing.T) {
	p := testAzureProvider("")
	p.SetGroupRestriction([]string{"test-group-2"})

	session := &sessions.SessionState{
		AccessToken: "imaginary_access_token",
		IDToken:     "imaginary_IDToken_token",
		Groups:      "no one|cares|non existing|groups|test-group-1"}
	result := p.ValidateGroup(session)

	assert.Equal(t, false, result)
}

func TestAzureRightPermittedGroups(t *testing.T) {
	p := testAzureProvider("")
	p.SetGroupRestriction([]string{"test-group-1", "test-group-2"})

	session := &sessions.SessionState{
		AccessToken: "imaginary_access_token",
		IDToken:     "imaginary_IDToken_token",
		Groups:      "no one|cares|test-group-2|non existing|groups"}
	result := p.ValidateGroup(session)

	assert.Equal(t, true, result)
}
*/
