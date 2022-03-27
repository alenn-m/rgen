package service_init

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"strings"

	"github.com/alenn-m/rgen/v2/generator/parser"
	"github.com/alenn-m/rgen/v2/util/config"
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

// ServiceInit generator
type ServiceInit struct {
	parsedData *parsedData
}

// Generate generates service_init
func (s *ServiceInit) Generate(input *parser.Parser, conf *config.Config) error {
	s.parseData(input, conf)

	f, err := s.getMainFileContent()
	if err != nil {
		return err
	}

	service := fmt.Sprintf(`"%s/api/%s"`, s.parsedData.Package, s.parsedData.LowerCaseName)
	if !strings.Contains(f, service) {
		authSvc := ", authSvc"
		if s.parsedData.Public {
			authSvc = ""
		}
		serviceInitToInsert := fmt.Sprintf(`// %s
			%s.New(r, db%s)`, inflection.Plural(s.parsedData.LowerCaseName), s.parsedData.LowerCaseName, authSvc)

		servicesIndex := strings.Index(f, "// [services]") + 13
		f = f[:servicesIndex] + fmt.Sprintf("\n%s\n", service) + f[servicesIndex:]

		if !s.parsedData.Public {
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

		s.parsedData.Content = string(content)
	}

	return nil
}

// Save saves generated service_init content to file
func (s *ServiceInit) Save() error {
	return ioutil.WriteFile(mainFileName, []byte(s.GetContent()), 0644)
}

// GetContent returns service_init generated content
func (s *ServiceInit) GetContent() string {
	return s.parsedData.Content
}

func (s *ServiceInit) parseData(input *parser.Parser, conf *config.Config) {
	s.parsedData = &parsedData{
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
