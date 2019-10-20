package view

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"

	"github.com/alecthomas/template"
	"github.com/leeif/pluto/config"
	plog "github.com/leeif/pluto/log"
)

var templates []string
var viewPath string

func Init(config *config.Config, logger *plog.PlutoLog) (err error) {
	templates = make([]string, 0)
	viewPath = config.View.Path
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
		if logger != nil {
			logger.Info(fmt.Sprintf("add template view: %s", filepath))
		}

		templates = append(templates, filepath)
	}

	return nil
}

func Parse(files ...string) (t *template.Template, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()
	paths := make([]string, 0)
	for _, file := range files {
		paths = append(paths, path.Join(viewPath, file))
	}
	tv := template.Must(template.ParseFiles(paths...))
	if err != nil {
		return nil, err
	}

	for _, tpl := range templates {
		tv = template.Must(tv.ParseFiles(tpl))
	}
	return tv, nil
}
