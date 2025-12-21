package wrapper

import (
	"embed"
	"fmt"
	"github.com/pm1381/sirish/internal"
	"github.com/pm1381/sirish/internal/dto"
	"github.com/pm1381/sirish/internal/visitors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

//go:embed test_samples/template/*.gotmpl
var f embed.FS

func TestAPMWrapperGenerator(t *testing.T) {
	type input struct {
		filename   string
		interfaces []string
	}
	type expected struct {
	}
	type scenario struct {
		name     string
		input    input
		expected expected
	}
	scenarios := []scenario{
		{
			name: "NoParamsNoResultTest",
			input: input{
				filename: "interface_samples.go",
				interfaces: []string{
					"NoParams",
					"NoResult",
				},
			},
			expected: expected{},
		},
		{
			name: "UnderscoreNamesTest",
			input: input{
				filename: "interface_samples.go",
				interfaces: []string{
					"UnderscoreNames",
				},
			},
			expected: expected{},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			path := internal.GetTestPathHelper(s.input.filename, "")
			var types = new(dto.Types)
			for _, ei := range s.input.interfaces {
				err := types.Set(ei)
				require.NoError(t, err)
			}
			typeVisitor := visitors.NewTypeVisitor(path, internal.GenerateUniqueValues(*types, nil))
			err := typeVisitor.Traverse()
			require.NoError(t, err)

			apmW := NewApmWrapper("sirish", "test_samples/template/wrapper.gotmpl", f, typeVisitor.GetWrappedInterfaces(), typeVisitor.GetImports())
			err = apmW.Generate(APMTypeWrapperOptions{
				GeneralOptions{
					Version:  "0.0.1",
					Imports:  true,
					CreateTx: true,
				},
			})
			require.NoError(t, err)
			sirishFiles := getFileNamesInDir(t, "test_samples")
			if len(s.input.interfaces) > 1 {
				// the naming structure differs
				var correctFileNames []string
				for _, eachIntr := range s.input.interfaces {
					correctFileNames = append(correctFileNames, fmt.Sprintf("%s.%s.sirish.go", eachIntr, strings.ReplaceAll(s.input.filename, ".go", "")))
				}
				fmt.Println(correctFileNames, sirishFiles)
				var foundCount int
				for _, wantName := range correctFileNames {
					for _, existFile := range sirishFiles {
						if wantName == existFile {
							foundCount++
						}
					}
				}
				assert.Equal(t, len(correctFileNames), foundCount)
			} else {
				var found bool
				generatedFile := fmt.Sprintf("%s.%s.go", strings.ReplaceAll(s.input.filename, ".go", ""), "sirish")
				for _, file := range sirishFiles {
					if file == generatedFile {
						found = true
					}
				}
				assert.True(t, found)
			}
		})
	}
	// removing generated files
	genFiles := getFileNamesInDir(t, "test_samples")
	for _, file := range genFiles {
		if strings.HasSuffix(file, ".sirish.go") {
			fmt.Println(file)
			// TODO: optional Remove
			//path := internal.GetTestPathHelper(file, "")
			//err := os.Remove(path)
			//require.NoError(t, err)
		}
	}
}

func getFileNamesInDir(t *testing.T, dir string) (res []string) {
	files, err := os.ReadDir(dir)
	require.NoError(t, err)
	for _, file := range files {
		if !file.IsDir() {
			if strings.Contains(file.Name(), ".sirish.go") {
				res = append(res, file.Name())
			}
		}
	}
	return res
}
