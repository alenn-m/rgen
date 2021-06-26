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
	LowerCaseName string
	Name          string
	Public        bool
}

type ServiceInit struct {
	Input  *Input
	Config *config.Config
}

func (s *ServiceInit) Init(input *Input, conf *config.Config) {
	s.Input = &Input{
		LowerCaseName: strings.ToLower(inflection.Singular(input.Name)),
		Name:          inflection.Singular(input.Name),
	}
	s.Config = conf
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

	service := fmt.Sprintf(`"%s/api/%s"`, s.Config.Package, s.Input.LowerCaseName)
	if !strings.Contains(f, service) {
		authSvc := ", authSvc"
		if s.Input.Public {
			authSvc = ""
		}
		serviceInitToInsert := fmt.Sprintf(`// %s
			%s.New(r, db%s)`, inflection.Plural(s.Input.LowerCaseName), s.Input.LowerCaseName, authSvc)

		servicesIndex := strings.Index(f, "// [services]") + 13
		f = f[:servicesIndex] + fmt.Sprintf("\n%s\n", service) + f[servicesIndex:]

		if !s.Input.Public {
			protectedRoutesIndex := strings.Index(f, "// [protected routes]") + 21
			f = f[:protectedRoutesIndex] + fmt.Sprintf("\n\n%s\n", serviceInitToInsert) + f[protectedRoutesIndex:]
		} else {
			publicRoutesIndex := strings.Index(f, "// [public routes]") + 18
			f = f[:publicRoutesIndex] + fmt.Sprintf("\n\n%s\n", serviceInitToInsert) + f[publicRoutesIndex:]
		}

		content, err := format.Source([]byte(f))
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(mainFilePath, content, 0644)

		return err
	}

	return nil
}
