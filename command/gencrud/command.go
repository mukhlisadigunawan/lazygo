package gencrud

import (
	"fmt"

	"github.com/mukhlisadigunawan/lazygo/command/gengateway"
	"github.com/mukhlisadigunawan/lazygo/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	EntityName  string
	DomainName  string
}

func Run(inputs ...string) error {

	packagePath := utils.GetPackagePath()

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create a new usecase\n" +
			"   gogen crud Product\n" +
			"     'Product' is an existing entity name\n" +
			"\n")

		return err
	}

	gcfg := utils.GetGogenConfig()

	entityName := inputs[0]

	obj := &ObjTemplate{
		PackagePath: packagePath,
		EntityName:  entityName,
		DomainName:  gcfg.Domain,
	}

	fileRenamer := map[string]string{
		"domainname": utils.LowerCase(gcfg.Domain),
		"entityname": utils.SnakeCase(entityName),
	}

	err := utils.CreateEverythingExactly("templates/", "shared", nil, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	// TODO check existing entity. if exist, read all the field else create new one

	err = utils.CreateEverythingExactly2(".gogen/templates/crud/", gcfg.Crud, fileRenamer, obj)
	if err != nil {
		return err
	}

	err = gengateway.Run(fmt.Sprintf("gateway%s", entityName))
	if err != nil {
		return err
	}

	return nil

}
