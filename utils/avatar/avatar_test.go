package avatar_test

import (
	"testing"

	"github.com/leeif/pluto/utils/avatar"
	"github.com/stretchr/testify/assert"
)

func TestRandomAvatar(t *testing.T) {
	avatar := avatar.Avatar{}
	url, _ := avatar.GetRandomAvatar()

	assert.NotEqual(t, url, "", "avatar url should not be empty")
}
