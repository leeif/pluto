package avatar_test

import (
	"io"
	"os"
	"testing"

	"github.com/leeif/pluto/utils/avatar"
)

func TestRandomAvatar(t *testing.T) {
	avatar := avatar.Avatar{}
	reader, err := avatar.GetRandomAvatar()
	if err != nil {
		t.Fatalf("expected to be no error, but err: %v", err)
	}
	file, _ := os.OpenFile("avatar.png", os.O_CREATE|os.O_RDWR, 0666)
	b, err := io.Copy(file, reader.Reader)
	defer func() {
		reader.Reader.Close()
		file.Close()
		os.Remove(file.Name())
	}()
	t.Logf("file extension: %v", reader.Ext)
	t.Logf("write file size: %v", b)
}
