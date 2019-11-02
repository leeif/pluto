package avatar_test

import (
	"testing"

	"github.com/leeif/pluto/utils/avatar"
	"github.com/stretchr/testify/assert"
)

func TestRandomAvatar(t *testing.T) {
	ag := avatar.AvatarGen{}
	ar, _ := ag.GenFromGravatar()
	url := ar.OriginURL

	assert.NotEqual(t, url, "", "avatar url should not be empty")
}
