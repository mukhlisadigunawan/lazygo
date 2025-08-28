package gengateway

import (
	"fmt"
	"lazygo/utils"
	"os"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	DomainName  string
	GatewayName string
	UsecaseName *string
	Methods     utils.OutportMethods
}

func Run(inputs ...string) error {

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create a gateway for all usecases with cloverdb sample implementation\n" +
			"   gogen gateway inmemory\n" +
			"     'inmemory' is a gateway name\n" +
			"\n")

		return err
	}

	packagePath := utils.GetPackagePath()
	gcfg := utils.GetGogenConfig()
	gatewayName := inputs[0]

	obj := ObjTemplate{
		PackagePath: packagePath,
		GatewayName: gatewayName,
		DomainName:  utils.LowerCase(gcfg.Domain),
		UsecaseName: nil,
	}

	driverName := gcfg.Gateway

	// first we create the shared
	// err := utils.CreateEverythingExactly("templates/", "shared", nil, obj, utils.AppTemplates)
	// if err != nil {
	// 	return err
	// }

	var notExistingMethod utils.OutportMethods

	// we read all the usecase folders
	//var folders []string
	fileInfo, err := os.ReadDir(fmt.Sprintf("domain_%s/usecase", gcfg.Domain))
	if err != nil {
		return err
	}

	uniqueMethodMap := map[string]int{}

	// trace all usecase
	for _, file := range fileInfo {

		// skip all the file
		if !file.IsDir() {
			continue
		}

		folderName := file.Name()

		em, err := createGatewayImpl(driverName, folderName, obj)
		if err != nil {
			return err
		}

		// we filter only the new method and skip the existing
		for _, method := range em {

			if _, exist := uniqueMethodMap[method.MethodName]; exist {
				continue
			}

			notExistingMethod = append(notExistingMethod, method)

			uniqueMethodMap[method.MethodName] = 1
		}
	}

	gatewayCode, err := getGatewayMethodTemplate(driverName)
	if err != nil {
		return err
	}

	// we will only inject the non existing method
	obj.Methods = notExistingMethod

	templateHasBeenInjected, err := utils.PrintTemplate(string(gatewayCode), obj)
	if err != nil {
		return err
	}

	gatewayFilename := fmt.Sprintf("domain_%s/gateway/%s/gateway.go", gcfg.Domain, gatewayName)

	bytes, err := injectToGateway(gatewayFilename, templateHasBeenInjected)
	if err != nil {
		return err
	}

	// reformat gateway.go
	err = utils.Reformat(gatewayFilename, bytes)
	if err != nil {
		return err
	}

	return nil

}

func createGatewayImpl(driverName, usecaseName string, obj ObjTemplate) (utils.OutportMethods, error) {
	outportMethods, err := utils.NewOutportMethods(obj.DomainName, usecaseName)
	if err != nil {
		return nil, err
	}

	obj.Methods = outportMethods
	err = utils.CreateEverythingExactly2(".lazygo/templates/gateway/", driverName, map[string]string{
		"gatewayname": utils.LowerCase(obj.GatewayName),
		"domainname":  obj.DomainName,
	}, obj)
	if err != nil {
		return nil, err
	}

	gatewayRootFolderName := fmt.Sprintf("domain_%s/gateway/%s", obj.DomainName, utils.LowerCase(obj.GatewayName))

	// file gateway impl file is already exist, we want to inject non existing method
	existingFunc, err := utils.NewOutportMethodImpl("gateway", gatewayRootFolderName)
	if err != nil {
		return nil, err
	}

	// collect the only methods that has not added yet
	notExistingMethod := utils.OutportMethods{}
	for _, m := range outportMethods {
		if _, exist := existingFunc[m.MethodName]; !exist {
			notExistingMethod = append(notExistingMethod, m)
		}
	}
	return notExistingMethod, nil
}

// getGatewayMethodTemplate ...
func getGatewayMethodTemplate(driverName string) ([]byte, error) {
	s := fmt.Sprintf(".lazygo/templates/gateway/%s/domain_${domainname}/gateway/${gatewayname}/~inject._go", driverName)
	return os.ReadFile(s)
}

func injectToGateway(gatewayFilename, injectedCode string) ([]byte, error) {
	return utils.InjectCodeAtTheEndOfFile(gatewayFilename, injectedCode)
}
