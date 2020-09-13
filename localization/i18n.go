package localization

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func NewBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("localization/active.zh.toml")
	bundle.MustLoadMessageFile("localization/active.en.toml")
	bundle.MustLoadMessageFile("localization/active.ja.toml")
	return bundle
}
