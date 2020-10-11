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
	s.Input = input
	s.Config = conf

	s.Input.LowerCaseName = strings.ToLower(inflection.Singular(s.Input.Name))
	s.Input.Name = inflection.Singular(s.Input.Name)
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
	repo := fmt.Sprintf(`%sDB "%s/api/%s/repositories/mysql"`, s.Input.LowerCaseName, s.Config.Package, s.Input.LowerCaseName)
	if !strings.Contains(f, service) && !strings.Contains(f, repo) {
		importToInsert := fmt.Sprintf("%s\n%s", service, repo)

		fmt.Println("service init public: ", s.Input.Public)

		authSvc := ", authSvc"
		if s.Input.Public {
			authSvc = ""
		}
		serviceInitToInsert := fmt.Sprintf(`// %s
			%s.New(r, %s.NewController(
				%sDB.New%sDB(db)%s),
			)`, inflection.Plural(s.Input.LowerCaseName), s.Input.LowerCaseName, s.Input.LowerCaseName, s.Input.LowerCaseName, strings.Title(s.Input.Name), authSvc)

		servicesIndex := strings.Index(f, "// [services]") + 13
		f = f[:servicesIndex] + fmt.Sprintf("\n%s\n", importToInsert) + f[servicesIndex:]

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
