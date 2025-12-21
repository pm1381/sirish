package visitors

import (
	"fmt"
	"github.com/pm1381/sirish/internal/dto"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path"
	"strconv"
)

type TypeVisitor struct {
	targetInterfaces  dto.Types
	fSet              *token.FileSet
	fileAbsPath       string
	packageName       string
	importAlias       dto.PkgImports
	wrappedInterfaces []dto.InterfaceInfo
	needMultipleFiles bool
}

func NewTypeVisitor(fileAbsPath string, targets dto.Types) *TypeVisitor {
	return &TypeVisitor{
		targetInterfaces:  targets,
		fileAbsPath:       fileAbsPath,
		importAlias:       make(map[string]string),
		wrappedInterfaces: make([]dto.InterfaceInfo, 0, len(targets)),
		needMultipleFiles: len(targets) > 1,
	}
}

func (tv *TypeVisitor) Traverse() error {
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, tv.fileAbsPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	tv.fSet = fSet
	ast.Walk(tv, file)
	return nil
}

func (tv *TypeVisitor) GetWrappedInterfaces() []dto.InterfaceInfo {
	return tv.wrappedInterfaces
}

func (tv *TypeVisitor) GetImports() dto.PkgImports {
	return tv.importAlias
}

func (tv *TypeVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch nodeWithType := node.(type) {
	//it does know the structure of arguments and types, just not their meaning.
	case *ast.File:
		tv.packageName = nodeWithType.Name.Name
	case *ast.ImportSpec:
		err := tv.handleImports(nodeWithType)
		if err != nil {
			log.Fatal(err)
		}
	case *ast.TypeSpec:
		switch interfaceType := nodeWithType.Type.(type) {
		case *ast.InterfaceType:
			if nodeWithType.Name == nil {
				return nil
			}
			interfaceName := nodeWithType.Name.Name
			if !tv.targetInterfaces.Exists(interfaceName) {
				return nil // means the interface is not in the list to search for
			}
			fn := path.Base(tv.fileAbsPath)
			if tv.needMultipleFiles {
				fn = interfaceName + "." + path.Base(tv.fileAbsPath) // something like name.profile_store.go
			}
			interfaceInfo := dto.InterfaceInfo{
				Name:      interfaceName,
				FilePath:  tv.fileAbsPath,
				FileName:  fn,
				Package:   tv.packageName,
				Directory: path.Dir(tv.fileAbsPath),
				Methods:   nil,
			}
			if nodeWithType.TypeParams != nil {
				err := tv.handleGenerics(nodeWithType.TypeParams, &interfaceInfo)
				if err != nil {
					log.Fatal(err)
				}
			}
			err := tv.handleInterface(interfaceType, &interfaceInfo)
			if err != nil {
				log.Fatal(err)
			}
			tv.wrappedInterfaces = append(tv.wrappedInterfaces, interfaceInfo)
			return nil // no need to check this interface children
		}
	}
	return tv
}

func (tv *TypeVisitor) handleInterface(node *ast.InterfaceType, interfaceDto *dto.InterfaceInfo) error {
	if node.Incomplete {
		return nil
	}
	var methods []dto.Method
	for _, method := range node.Methods.List {
		if len(method.Names) == 0 {
			continue // embed interfaces
		}
		switch functionWithType := method.Type.(type) {
		case *ast.FuncType:
			methodInfo := dto.Method{
				Name:        method.Names[0].Name,
				SpecialName: fmt.Sprintf("%s.%s", interfaceDto.Name, method.Names[0].Name),
				Params:      nil,
				Results:     nil,
				SpanName:    "span",
			}
			errParam := tv.handleParams(functionWithType.Params, &methodInfo)
			if errParam != nil {
				return errParam
			}
			errRes := tv.handleResults(functionWithType.Results, &methodInfo)
			if errRes != nil {
				return errRes
			}
			methods = append(methods, methodInfo)
		}
	}
	interfaceDto.Methods = methods
	return nil
}

func (tv *TypeVisitor) handleGenerics(genericFieldList *ast.FieldList, interfaceDto *dto.InterfaceInfo) error {
	if genericFieldList == nil || len(genericFieldList.List) == 0 {
		return nil
	}
	var typeParams []dto.TypeParamInfo
	for _, field := range genericFieldList.List {
		if field == nil || len(field.Names) == 0 {
			continue
		}
		constraint := ExprToString(tv.fSet, field.Type)
		if constraint == "" {
			constraint = "any"
		}
		for _, name := range field.Names {
			typeParams = append(typeParams, dto.TypeParamInfo{
				Name:       name.Name,
				Constraint: constraint,
			})
		}
	}
	interfaceDto.TypeParams = typeParams
	return nil
}

func (tv *TypeVisitor) handleParams(params *ast.FieldList, method *dto.Method) error {
	if params == nil || len(params.List) == 0 {
		return nil
	}
	var paramsInfo []dto.ParamInfo
	for index, p := range params.List {
		if p == nil {
			continue
		}
		typeStr := ExprToString(tv.fSet, p.Type)
		if !method.HasCtx && (typeStr == "context.Context") {
			method.HasCtx = true // handling first ctx occurs
			if len(p.Names) == 0 {
				ctxName := fmt.Sprintf("ctx_%d_0", index)
				paramsInfo = append(paramsInfo, dto.ParamInfo{
					Name: ctxName,
					Type: typeStr,
				})
				method.CtxName = ctxName
			} else {
				for i, name := range p.Names {
					if i == 0 {
						ctxName := fmt.Sprintf("ctx_%d_0", index)
						paramsInfo = append(paramsInfo, dto.ParamInfo{
							Name: ctxName,
							Type: typeStr,
						})
						method.CtxName = ctxName
					} else {
						if name.Name == "_" {
							paramsInfo = append(paramsInfo, dto.ParamInfo{
								Name: MethodParamSnowflake(method.Name, index, i, 4, "Un"),
								Type: typeStr,
							})
						} else {
							paramsInfo = append(paramsInfo, dto.ParamInfo{
								Name: name.Name,
								Type: typeStr,
							})
						}
					}
				}
			}
			continue
		}
		if len(p.Names) == 0 {
			paramsInfo = append(paramsInfo, dto.ParamInfo{
				Name: MethodParamSnowflake(method.Name, index, 0, 4, fmt.Sprintf("Un%s", typeStr)),
				Type: typeStr,
			})
		} else {
			for i, name := range p.Names {
				if name.Name == "_" {
					paramsInfo = append(paramsInfo, dto.ParamInfo{
						Name: MethodParamSnowflake(method.Name, index, i, 4, "Un"),
						Type: typeStr, // covering edge-cases like underscore(Us) params
					})
				} else {
					if name.Name == method.SpanName {
						method.SpanName = MethodParamSnowflake(method.Name, index, i, 4, "Spn")
					}
					paramsInfo = append(paramsInfo, dto.ParamInfo{
						Name: name.Name,
						Type: typeStr, // covering edge-cases like (a,b, c int)
					})
				}
			}
		}
	}

	var paramsNames string
	var ParamsOverallNames string
	for _, eachParam := range paramsInfo {
		paramsNames += eachParam.Name + ", "
		ParamsOverallNames += fmt.Sprintf("%s %s, ", eachParam.Name, eachParam.Type)
	}
	method.ParamsNames = paramsNames[:len(paramsNames)-2]
	method.ParamsOverallNames = ParamsOverallNames[:len(ParamsOverallNames)-2]

	method.Params = paramsInfo
	return nil
}

func (tv *TypeVisitor) handleResults(results *ast.FieldList, method *dto.Method) error {
	if results == nil || len(results.List) == 0 {
		return nil
	}
	var resultsInfo []dto.ResultInfo
	for index, p := range results.List {
		if p == nil {
			continue
		}
		typeStr := ExprToString(tv.fSet, p.Type)
		if typeStr == "error" {
			method.HasError = true
		}
		if len(p.Names) == 0 {
			resultsInfo = append(resultsInfo, dto.ResultInfo{
				Name: MethodParamSnowflake(method.Name, index, 0, 4, "ResUn"),
				Type: typeStr,
			})
		} else {
			// not really casual :)
			method.HasNamedResult = true
			for i, name := range p.Names {
				if name.Name == method.SpanName {
					method.SpanName = MethodParamSnowflake(method.Name, index, i, 4, "Spn")
				}
				n := name.Name
				if name.Name == "_" {
					n = MethodParamSnowflake(method.Name, index, i, 4, "ResUn")
				}
				resultsInfo = append(resultsInfo, dto.ResultInfo{
					Name: n,
					Type: typeStr,
				})
			}
		}
	}
	method.Results = resultsInfo

	var resultNames string
	var resultTypesNames string
	var resultOverallNames string
	for _, eachResult := range resultsInfo {
		resultNames += eachResult.Name + ", "
		resultTypesNames += eachResult.Type + ", "
		resultOverallNames += fmt.Sprintf("%s %s, ", eachResult.Name, eachResult.Type)
	}
	method.ResultNames = resultNames[:len(resultNames)-2]
	method.ResultTypesNames = resultTypesNames[:len(resultTypesNames)-2]
	method.ResultOverallNames = resultOverallNames[:len(resultOverallNames)-2]
	return nil
}

func (tv *TypeVisitor) handleImports(node *ast.ImportSpec) error {
	unquoteImport, err := strconv.Unquote(node.Path.Value)
	if err != nil {
		return err
	}
	var alias string
	if node.Name != nil {
		alias = node.Name.Name
	} else {
		alias = path.Base(unquoteImport)
	}
	tv.importAlias[unquoteImport] = alias
	return nil
}
