package dto

// Method defines how wrapper method signature will be
type Method struct {
	Name               string
	SpecialName        string
	Params             []ParamInfo
	Results            []ResultInfo
	HasNamedResult     bool
	HasCtx             bool
	ResultNames        string
	ResultTypesNames   string
	ResultOverallNames string
	ParamsNames        string
	ParamsOverallNames string
	HasError           bool
	ErrorName          string
	CtxName            string
	SpanName           string
}

type TypeParamInfo struct {
	Name       string // for example T or K
	Constraint string // for example any, comparable or {interface | int64}
}

type ParamInfo struct {
	Name string // can be empty
	Type string // printable type, e.g. "context.Context"
}

type ResultInfo struct {
	Name string // can be empty
	Type string // printable type, e.g. "context.Context"
}

type InterfaceInfo struct {
	Name       string
	TypeParams []TypeParamInfo
	Methods    []Method
	FileName   string
	FilePath   string
	Package    string
	Directory  string
}

// PkgImports is a map of imports which their key is path and value is possible alias
type PkgImports map[string]string

func (pi PkgImports) PathExists(path string) bool {
	for key, _ := range pi {
		if key == path {
			return true
		}
	}
	return false
}
