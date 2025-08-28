package genusecase

import (
	"fmt"
	"lazygo/utils"
	"strings"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	UsecaseName string
}

func Run(inputs ...string) error {
	var (
		err error
	)
	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create a new usecase\n" +
			"   lazygo usecase {CamelCaseName}\n" +
			"   better use ResourceNameMethod format \n" +
			"   ex: MemberCreate, MemberDelete, MemberGetOne, MemberGetAll, MemberUpdate\n")

		return err
	}

	packagePath := utils.GetPackagePath()
	gcfg := utils.GetGogenConfig()

	usecaseName := inputs[0]

	obj := &ObjTemplate{
		PackagePath: packagePath,
		UsecaseName: usecaseName,
	}

	fileRenamer := map[string]string{
		"usecasename": utils.SnakeCase(usecaseName),
		"domainname":  gcfg.Domain,
	}

	if strings.HasSuffix(utils.LowerCase(usecaseName), "getall") {
		err = utils.CreateEverythingExactly("templates/usecase/", "getall", fileRenamer, obj, utils.AppTemplates)
		if err != nil {
			return err
		}
	} else if strings.HasPrefix(utils.LowerCase(usecaseName), "run") {
		err = utils.CreateEverythingExactly("templates/usecase/", "run", fileRenamer, obj, utils.AppTemplates)
		if err != nil {
			return err
		}
	} else {
		err = utils.CreateEverythingExactly("templates/usecase/", "default", fileRenamer, obj, utils.AppTemplates)
		if err != nil {
			return err
		}
	}

	return nil

}
