package view

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"

	"github.com/alecthomas/template"
	"github.com/MuShare/pluto/config"
	"golang.org/x/text/language"
)

var view *View

type View struct {
	templates        []string
	path             string
	supportLanguages []language.Tag
}

func GetView() (*View, error) {
	if view != nil {
		return view, nil
	}

	return nil, errors.New("view not init")
}

func InitView(config *config.Config) (err error) {
	templates := make([]string, 0)
	templatePath := path.Join(config.View.Path, "template")
	files, err := ioutil.ReadDir(templatePath)
	if err != nil {
		return err
	}
	rep, err := regexp.Compile(`\.html$`)
	if err != nil {
		return err
	}
	for _, file := range files {
		filepath := path.Join(templatePath, file.Name())
		if !rep.MatchString(filepath) {
			continue
		}
		log.Println(fmt.Sprintf("add template view: %s", filepath))

		templates = append(templates, filepath)
	}

	// support language
	supported, _, err := language.ParseAcceptLanguage(config.View.Languages)

	if err != nil {
		return err
	}

	view = &View{
		path:             config.View.Path,
		templates:        templates,
		supportLanguages: supported,
	}

	return nil
}

func (view *View) IsExist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func (v *View) GetMatchLanguage(accpetLanguageHeader string) (string, error) {
	tags, _, err := language.ParseAcceptLanguage(accpetLanguageHeader)

	if err != nil {
		return "", err
	}

	languageMatcher := language.NewMatcher(v.supportLanguages)

	tag, _, _ := languageMatcher.Match(tags...)
	base, _ := tag.Base()
	return base.String(), nil
}

func (v *View) GetMatchLanguageFile(file string, lan string) (string, error) {
	lanFile := path.Join(v.path, path.Join(lan, file))
	if v.IsExist(lanFile) {
		return lanFile, nil
	}

	for _, tag := range v.supportLanguages {
		base, _ := tag.Base()
		lanFile := path.Join(v.path, path.Join(base.String(), file))
		if v.IsExist(lanFile) {
			return lanFile, nil
		}
	}

	return "", errors.New("file not exists")
}

func (v *View) Parse(acceptLanguageHeader string, files ...string) (t *template.Template, err error) {

	lan, err := v.GetMatchLanguage(acceptLanguageHeader)

	if err != nil {
		return nil, err
	}

	paths := make([]string, 0)
	for _, file := range files {
		lanFile, err := v.GetMatchLanguageFile(file, lan)
		if err != nil {
			return nil, err
		}
		paths = append(paths, lanFile)
	}

	tv, err := template.ParseFiles(paths...)
	if err != nil {
		return nil, err
	}

	for _, tpl := range v.templates {
		tv, err = tv.ParseFiles(tpl)
		if err != nil {
			return nil, err
		}
	}

	return tv, nil
}
