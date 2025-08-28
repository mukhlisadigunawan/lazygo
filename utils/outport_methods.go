package utils

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// OutportMethods ...
type OutportMethods []*method

type method struct {
	InterfaceName    string //
	MethodName       string //
	MethodSignature  string //
	DefaultParamVal  string //
	DefaultReturnVal string //
}

// NewOutportMethods read the Outport interface and collect all method on it
func NewOutportMethods(domainName, usecaseName string) (OutportMethods, error) {
	fileReadPath := fmt.Sprintf("domain_%s/usecase/%s", domainName, strings.ToLower(usecaseName))

	var om OutportMethods

	err := om.readInterface("", "Outport", fileReadPath, GetPackagePath())
	if err != nil {
		return nil, err
	}

	return om, nil
}

// readInterface this method is used recursively
func (obj *OutportMethods) readInterface(packageName, interfaceName, folderPath, packagePath string) error {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, folderPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {

		// read file by file
		for _, file := range pkg.Files {

			ip := map[string]string{}

			// for each file, we will read line by line
			for _, decl := range file.Decls {

				gen, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				for _, spec := range gen.Specs {

					is, ok := spec.(*ast.ImportSpec)
					if ok {
						handleImports(packagePath, is, ip)
						continue
					}

					// Outport is must a type spec
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					// start by looking the Outport interface with name "Outport"
					if ts.Name.String() != interfaceName {
						continue
					}

					// make sure Outport is an interface
					iFace, ok := ts.Type.(*ast.InterfaceType)
					if !ok {
						continue
					}

					for _, field := range iFace.Methods.List {

						// as a field, there are two possibility
						switch ty := field.Type.(type) {

						case *ast.SelectorExpr: // as extension of another interface
							expression := ty.X.(*ast.Ident).String()
							pathWithGomod := ip[expression]
							pathOnly := strings.TrimPrefix(pathWithGomod, packagePath+"/")
							interfaceName := ty.Sel.String()
							//fmt.Printf(">>>>> %v %v\n", expression, interfaceName)

							err := obj.readInterface(fmt.Sprintf("%s.%s", expression, interfaceName), interfaceName, pathOnly, packagePath)
							if err != nil {
								return err
							}

						case *ast.FuncType: // as direct func (method) interface
							//TODO cannot handle Something(c context.Context, a Hoho) yet, where the Hoho part is a struct
							err := obj.handleMethodSignature(packageName, file.Name.String(), ty, field.Names[0].String())
							if err != nil {
								return err
							}

						case *ast.Ident: // as interface extension in same package
							//ast.Print(fset, ty)
							//TODO as interface extension in same package in the same or different file
							fmt.Printf("as interface extension in same package in the same or different file not supported yet\n")
						}

					}

				}

			}

		}
	}

	return nil
}

func (obj *OutportMethods) handleMethodSignature(packageAndInterfaceName string, prefixExpression string, fType *ast.FuncType, methodName string) error {

	// checking first params as context.Context
	if !obj.validateFirstParamIsContext(fType) {
		return fmt.Errorf("function `%s` must have context.Context in its first param argument", methodName)
	}

	if fType.Results == nil {
		return fmt.Errorf("function `%s` result at least have error return value", methodName)
	}

	defParVal := obj.composeDefaultValue(fType.Params.List)
	defRetVal := obj.composeDefaultValue(fType.Results.List)

	methodSignature := TypeHandler{PrefixExpression: prefixExpression}.processFuncType(&bytes.Buffer{}, fType)
	msObj := method{
		InterfaceName:    packageAndInterfaceName,
		MethodName:       methodName,
		MethodSignature:  methodSignature,
		DefaultParamVal:  defParVal,
		DefaultReturnVal: defRetVal,
	}

	*obj = append(*obj, &msObj)

	return nil
}

func (obj *OutportMethods) composeDefaultValue(list []*ast.Field) string {
	defRetVal := ""
	for i, retList := range list {

		var v string

		switch t := retList.Type.(type) {

		case *ast.SelectorExpr:
			v = fmt.Sprintf("%v.%v{}", t.X, t.Sel)

			// little hacky for context
			if v == "context.Context{}" {
				v = "ctx"
			}

		case *ast.StarExpr:
			v = "nil"

		case *ast.Ident:

			if t.Name == "error" {
				v = "nil"

			} else if strings.HasPrefix(t.Name, "int") {
				v = "0"

			} else if strings.HasPrefix(t.Name, "float") {
				v = "0.0"

			} else if t.Name == "string" {
				v = "\"\""

			} else if t.Name == "bool" {
				v = "false"

			} else {
				v = "nil"
			}

		default:
			v = "nil"

		}

		// append the comma
		if i < len(list)-1 {
			defRetVal += v + ", "
		} else {
			defRetVal += v
		}

	}
	return defRetVal
}

func (obj *OutportMethods) validateFirstParamIsContext(fType *ast.FuncType) bool {

	if fType.Params == nil || len(fType.Params.List) == 0 {
		return false
	}

	se, ok := fType.Params.List[0].Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	firstParamArgs := fmt.Sprintf("%s.%s", se.X.(*ast.Ident).String(), se.Sel.String())
	if firstParamArgs != "context.Context" {
		return false
	}

	return true
}

func handleImports(packagePath string, is *ast.ImportSpec, ip map[string]string) {

	//fmt.Printf(">>> %s === %s : %v \n", fmt.Sprintf("\"%s", packagePath), is.Path.Value, strings.HasPrefix(is.Path.Value, fmt.Sprintf("\"%s", packagePath)))

	if !strings.HasPrefix(is.Path.Value, fmt.Sprintf("\"%s", packagePath)) {
		return
	}

	v := strings.Trim(is.Path.Value, "\"")
	if is.Name != nil {
		ip[is.Name.String()] = v
	} else {
		ip[v[strings.LastIndex(v, "/")+1:]] = v
	}
}
