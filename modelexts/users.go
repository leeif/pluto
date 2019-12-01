package modelexts

import (
	"github.com/leeif/pluto/models"
)

type User struct {
	*models.User
	Roles []string `json:"roles"`
}
