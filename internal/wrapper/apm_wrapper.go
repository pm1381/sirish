package wrapper

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"github.com/pm1381/sirish/internal/dto"
	"golang.org/x/tools/imports"
	"os"
	"path"
	"strings"
	"text/template"
)

type apmWrapper struct {
	suffix     string
	template   *template.Template
	interfaces []dto.InterfaceInfo
	imports    dto.PkgImports
}

type generatorValues struct {
	Version   string
	Interface dto.InterfaceInfo
	Imports   dto.PkgImports
	Suffix    string
	TypeName  string
	CreateTx  bool
}

const APMPath = "go.elastic.co/apm/v2"

func NewApmWrapper(suffix string, pattern string, f embed.FS, interfaces []dto.InterfaceInfo, imports dto.PkgImports) WrapperInterface {
	if suffix == "" {
		suffix = "sirish"
	}
	if pattern == "" {
		pattern = "internal/templates/wrapper.gotmpl"
	}
	if !imports.PathExists(APMPath) {
		imports[APMPath] = "apm" // add APM paths
	}
	tw := &apmWrapper{
		suffix:     suffix,
		interfaces: interfaces,
		imports:    imports,
		template:   template.Must(template.ParseFS(f, pattern)),
	}
	return tw
}

func (tw *apmWrapper) Generate(opts Options) error {
	options, ok := opts.(APMTypeWrapperOptions)
	if !ok {
		return errors.New("invalid options")
	}
	for _, eachInterface := range tw.interfaces {
		fmt.Printf("- generating sirish for interface %s in directory %s \n", eachInterface.Name, eachInterface.Directory)

		var err error
		buf := new(bytes.Buffer)
		filenameSuffix := fmt.Sprintf("%s.%s.go", strings.ReplaceAll(eachInterface.FileName, ".go", ""), tw.suffix) // profile_store.go ---> profile_store.sirish.go  OR interfaceName.profile_store.go ---> interfaceName.profile_store.sirish.go
		fullPath := path.Join(eachInterface.Directory, filenameSuffix)
		var processed []byte

		if err = tw.template.ExecuteTemplate(buf, "wrapper.gotmpl", generatorValues{
			Version:   options.Version,
			Interface: eachInterface,
			Imports:   tw.imports,
			Suffix:    tw.suffix,
			CreateTx:  options.CreateTx,
			TypeName:  eachInterface.Name + strings.Replace(tw.suffix, string(tw.suffix[0]), strings.ToUpper(string(tw.suffix[0])), 1),
		}); err != nil {
			fmt.Printf("error occured on execution: %v", err)
			continue
		}
		if options.Imports {
			processed, err = tw.formatImports(fullPath, buf)
			if err != nil {
				fmt.Printf("error formatting imports: %v", err)
				processed = buf.Bytes()
			}
		} else {
			processed = buf.Bytes()
		}
		f, err := os.Create(fullPath)
		if err != nil {
			fmt.Printf("error creating file: %v", err)
			continue
		}
		if _, err := f.Write(processed); err != nil {
			fmt.Printf("error writing to file: %v", err)
			f.Close()
			continue
		}
		f.Close()
	}
	return nil
}

func (tw *apmWrapper) formatImports(fileAbsPath string, buffer *bytes.Buffer) ([]byte, error) {
	formatedFile, err := imports.Process(fileAbsPath, buffer.Bytes(), nil)
	if err != nil {
		return nil, err
	}
	return formatedFile, nil
}
