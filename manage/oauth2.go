package manage

import (
	"github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"net/url"
)

func (m *Manager) AuthorizationCodeGrant(ot *request.OAuth2Tokens) {

}

func (m *Manager) PasswordGrant(ot *request.OAuth2Tokens) {

}

func (m *Manager) ClientCredentialGrant(ot *request.OAuth2Tokens) {

}

func (m *Manager) RefreshTokenGrant(ot *request.OAuth2Tokens) {

}

func (m *Manager) GrantAuthorizationCode(oa *request.OAuth2Authorize) (*url.URL, *pluto_error.PlutoError) {
	return nil, nil
}

func (m *Manager) GrantAccessToken(oa *request.OAuth2Authorize) (*url.URL, *pluto_error.PlutoError) {
	return nil, nil
}
