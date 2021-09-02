package providers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	//"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/apis/sessions"
)

var (
	path_group            string = "/v1.0/me/memberOf?$select=displayName"
	path_group_next       string = "/v1.0/me/memberOf?$select=displayName&$skiptoken=X%27test-token%27"
	path_group_wrong      string = "/v1.0/him/memberOf?$select=displayName"
	payload_group_empty   string = `{"@odata.context":"https://graph.microsoft.com/v1.0/$metadata#directoryObjects(displayName)","value":[]}`
	payload_group_garbage string = `{"@odata.context":"https://graph.microsoft.com/v1.0/$metadata#directoryObjects(displayName)","value":[{"@odata.type":"#microsoft.graph.group","displayName":"test-group-1"},{"@odata.type":"#microsoft.graph.group","displayName":"test-group-2"}]}`
	payload_group_simple  string = `{"@odata.context":"https://graph.microsoft.com/v1.0/$metadata#directoryObjects(displayName)","value":[{"@odata.type":"#microsoft.graph.group","displayName":"test-group-1"},{"@odata.type":"#microsoft.graph.group","displayName":"test-group-2"}]}`
	payload_group_part_1  string = `{"@odata.context":"https://graph.microsoft.com/v1.0/$metadata#directoryObjects(displayName)","@odata.nextLink":"https://graph.microsoft.com/v1.0/me/memberOf?$select=displayName&$skiptoken=X%27test-token%27","value":[{"@odata.type":"#microsoft.graph.group","displayName":"test-group-1"},{"@odata.type":"#microsoft.graph.group","displayName":"test-group-2"}]}`
	payload_group_part_2  string = `{"@odata.context":"https://graph.microsoft.com/v1.0/$metadata#directoryObjects(displayName)","value":[{"@odata.type":"#microsoft.graph.group","displayName":"test-group-3"}]}`
)

type mockTransport struct {
	params map[string]string
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

func newMockTransport(params map[string]string) http.RoundTripper {
	return &mockTransport{params}
}

func TestAzureProviderNoGroups(t *testing.T) {
	params := map[string]string{
		path_group: payload_group_empty}

	http.DefaultClient.Transport = newMockTransport(params)

	p := testAzureProvider("")

	session := &sessions.SessionState{
		AccessToken: "imaginary_access_token",
		IDToken:     "imaginary_IDToken_token"}

	groups, err := p.GetGroups(session, "")
	http.DefaultClient.Transport = nil

	assert.Equal(t, nil, err)
	assert.Equal(t, "", groups)
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

/*ValidateGroup is implemented as general  SetAllowedGroups
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
func TestAzureLoginURLnoResource(t *testing.T) {
	p := testAzureProvider("")
	p.ProtectedResource = nil

	result := p.GetLoginURL("http://redirect/url", "state", "")
	params, _ := url.ParseQuery(result)
	nonce := params.Get("nonce")
	expecteURL := ""
	if nonce != "" {
		expecteURL = "?client_id=&nonce=" + nonce + "&prompt=&redirect_uri=http%3A%2F%2Fredirect%2Furl&response_mode=form_post&response_type=id_token+code&scope=openid&state=state"
	}

	assert.Equal(t, expecteURL, result)
}

func TestAzureLoginURL(t *testing.T) {
	p := testAzureProvider("")

	result := p.GetLoginURL("http://redirect/url", "state", "")
	params, _ := url.ParseQuery(result)
	nonce := params.Get("nonce")
	expecteURL := ""
	if nonce != "" {
		expecteURL = "?client_id=&nonce=" + nonce + "&prompt=&redirect_uri=http%3A%2F%2Fredirect%2Furl&resource=https%3A%2F%2Fgraph.microsoft.com&response_mode=form_post&response_type=id_token+code&scope=openid&state=state"
	}

	assert.Equal(t, expecteURL, result)
}
