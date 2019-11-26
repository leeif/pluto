package manage

import (
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
)

func (m *Manager) CreateRole(cr request.CreateRole) *perror.PlutoError {
	return nil
}

func (m *Manager) CreateScope(cs request.CreateScope) *perror.PlutoError {
	return nil
}

func (m *Manager) CreateApplication(ca request.CreateApplication) *perror.PlutoError {
	return nil
}

func (m *Manager) AttachScope(as request.RoleScope) *perror.PlutoError {
	return nil
}

func (m *Manager) DetachScope(as request.RoleScope) *perror.PlutoError {
	return nil
}

func (m *Manager) ApplicationDefaultRole(ar request.ApplicationRole) *perror.PlutoError {
	return nil
}
