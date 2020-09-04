package view_test

import (
	"testing"

	"github.com/MuShare/pluto/config"
	"github.com/MuShare/pluto/utils/view"
	"github.com/stretchr/testify/assert"
)

func Init() error {
	config, _ := config.NewConfig([]string{}, "test")
	config.View.Path = "../../views"
	if err := view.InitView(config); err != nil {
		return err
	}

	return nil
}

func TestGetMatchLanguage(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Expect no error, but: %v", err)
	}
	vw, _ := view.GetView()
	lan1, err := vw.GetMatchLanguage("zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7")
	if err != nil {
		t.Fatalf("Expect no error, but: %v", err)
	}

	assert.Equal(t, "zh", lan1, "language should be zh")

	lan2, err := vw.GetMatchLanguage("ja;q=0.7")
	if err != nil {
		t.Fatalf("Expect no error, but: %v", err)
	}

	assert.Equal(t, "ja", lan2, "language should be ja")

	lan3, err := vw.GetMatchLanguage("fr;q=1")
	if err != nil {
		t.Fatalf("Expect no error, but: %v", err)
	}

	assert.Equal(t, "en", lan3, "language should be en")

	lan4, err := vw.GetMatchLanguage("")
	if err != nil {
		t.Fatalf("Expect no error, but: %v", err)
	}

	assert.Equal(t, "en", lan4, "language should be en")
}

func TestGetMatchLanguageFile(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Expect no error, but: %v", err)
	}

	vw, _ := view.GetView()
	if _, err := vw.GetMatchLanguageFile("404.html", "en"); err != nil {
		t.Fatalf("Expect no error, but: %v", err)
	}

	if _, err := vw.GetMatchLanguageFile("404.html", "zh"); err != nil {
		t.Fatalf("Expect no error, but: %v", err)
	}

	if _, err := vw.GetMatchLanguageFile("404.html", "ja"); err != nil {
		t.Fatalf("Expect no error, but: %v", err)
	}

}
