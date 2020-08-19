package service_init

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/alenn-m/rgen/util/config"
	"github.com/jinzhu/inflection"
)

type Input struct {
	Name string
}

type ServiceInit struct {
	Input  *Input
	Config *config.Config
}

func (s *ServiceInit) Init(input *Input, conf *config.Config) {
	s.Input = input
	s.Config = conf

	s.Input.Name = strings.ToLower(inflection.Singular(s.Input.Name))
}

func (s *ServiceInit) Generate() error {
	mainFilePath, err := filepath.Abs("main.go")
	if err != nil {
		return err
	}

	mainFile, err := ioutil.ReadFile(mainFilePath)
	if err != nil {
		return err
	}

	f := string(mainFile)

	service := fmt.Sprintf(`"%s/api/%s"`, s.Config.Package, s.Input.Name)
	repo := fmt.Sprintf(`%sDB "%s/api/%s/repositories/mysql"`, s.Input.Name, s.Config.Package, s.Input.Name)
	if !strings.Contains(f, service) && !strings.Contains(f, repo) {
		importToInsert := fmt.Sprintf("%s\n%s", service, repo)

		serviceInitToInsert := fmt.Sprintf(`// %s
			%s.New(r, %s.NewController(
				%sDB.New%sDB(db), authSvc),
			)`, inflection.Plural(s.Input.Name), s.Input.Name, s.Input.Name, s.Input.Name, strings.Title(s.Input.Name))

		servicesIndex := strings.Index(f, "// [services]") + 13
		f = f[:servicesIndex] + fmt.Sprintf("\n%s\n", importToInsert) + f[servicesIndex:]

		publicRoutesIndex := strings.Index(f, "// [protected routes]") + 21
		f = f[:publicRoutesIndex] + fmt.Sprintf("\n\n%s\n", serviceInitToInsert) + f[publicRoutesIndex:]

		content, err := format.Source([]byte(f))
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(mainFilePath, content, 0644)

		return err
	}

	return nil
}
