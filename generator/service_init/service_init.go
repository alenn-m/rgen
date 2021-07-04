package service_init

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"strings"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/jinzhu/inflection"
)

var mainFileName = "main.go"

type parsedData struct {
	LowerCaseName string
	Name          string
	Public        bool
	Package       string
	Content       string
}

type ServiceInit struct {
	Input *parsedData
}

func (s *ServiceInit) Generate(input *parser.Parser, conf *config.Config) error {
	s.parseData(input, conf)

	f, err := s.getMainFileContent()
	if err != nil {
		return err
	}

	service := fmt.Sprintf(`"%s/api/%s"`, s.Input.Package, s.Input.LowerCaseName)
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

		s.Input.Content = string(content)
	}

	return nil
}

func (s *ServiceInit) Save() error {
	return ioutil.WriteFile(mainFileName, []byte(s.GetContent()), 0644)
}

func (s *ServiceInit) GetContent() string {
	return s.Input.Content
}

func (s *ServiceInit) parseData(input *parser.Parser, conf *config.Config) {
	s.Input = &parsedData{
		LowerCaseName: strings.ToLower(inflection.Singular(input.Name)),
		Name:          inflection.Singular(input.Name),
		Package:       conf.Package,
		Public:        input.Public,
	}
}

func (s *ServiceInit) getMainFileContent() (string, error) {
	f, err := ioutil.ReadFile(mainFileName)
	if err != nil {
		return "", err
	}

	return string(f), nil
}

func (s *ServiceInit) setMainFileLocation(location string) {
	mainFileName = location
}
