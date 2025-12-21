package visitors

import (
	"github.com/pm1381/sirish/internal"
	"github.com/pm1381/sirish/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"reflect"
	"testing"
)

func TestTypeVisitorWorks(t *testing.T) {
	type input struct {
		filename   string
		interfaces []string
	}
	type methodDetails struct {
		paramCount      int
		resultCount     int
		hasCtx          bool
		hasError        bool
		hasNamedResult  bool
		specialName     string
		name            string
		resultTypeNames string
	}
	type result struct {
		interfaceName   string
		methodSignature []methodDetails
	}
	type scenario struct {
		name   string
		input  input
		result result
	}

	scenarios := []scenario{
		{
			name: "NamedParamsAndResultsTest",
			input: input{
				filename:   "types.go",
				interfaces: []string{"NamedParamsAndResults"},
			},
			result: result{
				interfaceName: "NamedParamsAndResults",
				methodSignature: []methodDetails{
					{
						// Method1(a string, b *int, c []byte) (s string, err error)
						name:            "Method1",
						specialName:     "NamedParamsAndResults.Method1",
						paramCount:      3,
						resultCount:     2,
						hasCtx:          false,
						hasError:        true,
						hasNamedResult:  true,
						resultTypeNames: "string, error",
					},
				},
			},
		},
		{
			name: "NoParamsTest",
			input: input{
				filename:   "types.go",
				interfaces: []string{"NoParams"},
			},
			result: result{
				interfaceName: "NoParams",
				methodSignature: []methodDetails{
					{
						//Method1() error
						name:            "Method1",
						specialName:     "NoParams.Method1",
						paramCount:      0,
						resultCount:     1,
						hasCtx:          false,
						hasError:        true,
						hasNamedResult:  false,
						resultTypeNames: "error",
					},
					{
						//Method2() (s string, err error)
						name:            "Method2",
						specialName:     "NoParams.Method2",
						paramCount:      0,
						resultCount:     2,
						hasCtx:          false,
						hasError:        true,
						hasNamedResult:  true,
						resultTypeNames: "string, error",
					},
					{
						//Method3() (string, error)
						name:            "Method3",
						specialName:     "NoParams.Method3",
						paramCount:      0,
						resultCount:     2,
						hasCtx:          false,
						hasError:        true,
						hasNamedResult:  false,
						resultTypeNames: "string, error",
					},
				},
			},
		},
		{
			name: "NoResultTest",
			input: input{
				filename:   "types.go",
				interfaces: []string{"NoResult"},
			},
			result: result{
				interfaceName: "NoResult",
				methodSignature: []methodDetails{
					{
						//Method1(s []string)
						name:            "Method1",
						specialName:     "NoResult.Method1",
						paramCount:      1,
						resultCount:     0,
						hasCtx:          false,
						hasError:        false,
						hasNamedResult:  false,
						resultTypeNames: "",
					},
					{
						//Method2(a, b, c int, s context.Context)
						name:            "Method2",
						specialName:     "NoResult.Method2",
						paramCount:      4,
						resultCount:     0,
						hasCtx:          true,
						hasError:        false,
						hasNamedResult:  false,
						resultTypeNames: "",
					},
					{
						//Method3(a, _, c int, _ *context.Context)
						name:            "Method3",
						specialName:     "NoResult.Method3",
						paramCount:      4,
						resultCount:     0,
						hasCtx:          false,
						hasError:        false,
						hasNamedResult:  false,
						resultTypeNames: "",
					},
				},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			abs := internal.GetTestPathHelper(s.input.filename, "visitors")
			dt := new(dto.Types)
			for _, intr := range s.input.interfaces {
				err := dt.Set(intr)
				require.NoError(t, err)
			}
			wrapper := NewTypeVisitor(abs, internal.GenerateUniqueValues(*dt, nil))
			err := wrapper.Traverse()
			require.NoError(t, err)

			// assert
			interfaces := wrapper.GetWrappedInterfaces()
			for _, eachInterface := range interfaces {
				assert.Equal(t, s.result.interfaceName, eachInterface.Name)
				assert.Equal(t, len(s.result.methodSignature), len(eachInterface.Methods))
				for _, eachMethod := range eachInterface.Methods {
					var found bool
					for _, det := range s.result.methodSignature {
						if det.name == eachMethod.Name {
							found = true
							assert.Equal(t, det.name, eachMethod.Name)
							assert.Equal(t, det.specialName, eachMethod.SpecialName)
							assert.Equal(t, det.paramCount, len(eachMethod.Params))
							assert.Equal(t, det.resultCount, len(eachMethod.Results))
							assert.Equal(t, det.resultTypeNames, eachMethod.ResultTypesNames)
							assert.Equal(t, det.hasError, eachMethod.HasError)
							assert.Equal(t, det.hasCtx, eachMethod.HasCtx)
							assert.Equal(t, det.hasNamedResult, eachMethod.HasNamedResult)
						}
					}
					assert.True(t, found)
				}
			}
		})
	}
}

func TestPackageNameWorks(t *testing.T) {
	type input struct {
		filename   string
		interfaces []string
	}
	type result struct {
		pkgName       string
		interfaceName string
		fileName      string
	}
	type scenario struct {
		name   string
		input  input
		result result
	}
	scenarios := []scenario{
		{
			name: "pkgNameTest",
			input: input{
				filename:   "imports.go",
				interfaces: []string{"RandConflict"},
			},
			result: result{
				pkgName:       "test_samples",
				interfaceName: "RandConflict",
				fileName:      "imports.go",
			},
		},
	}
	for _, each := range scenarios {
		t.Run(each.name, func(t *testing.T) {
			abs := internal.GetTestPathHelper(each.input.filename, "visitors")
			dt := new(dto.Types)
			for _, intr := range each.input.interfaces {
				err := dt.Set(intr)
				require.NoError(t, err)
			}
			wrapper := NewTypeVisitor(abs, internal.GenerateUniqueValues(*dt, nil))
			err := wrapper.Traverse()
			require.NoError(t, err)

			// assert
			interfaces := wrapper.GetWrappedInterfaces()
			for _, eachInterface := range interfaces {
				assert.Equal(t, each.result.pkgName, eachInterface.Package)
				assert.Equal(t, each.result.fileName, eachInterface.FileName)
				assert.Equal(t, each.result.interfaceName, eachInterface.Name)
				assert.Equal(t, abs, eachInterface.FilePath)

				fullPath := filepath.Join(eachInterface.Directory, eachInterface.FileName)
				assert.Equal(t, fullPath, eachInterface.FilePath)
			}
		})
	}
}

func TestImportsWorks(t *testing.T) {
	type input struct {
		filename   string
		interfaces []string
	}
	type result struct {
		expectedImports dto.PkgImports
	}
	type scenario struct {
		name   string
		input  input
		result result
	}
	scenarios := []scenario{
		{
			name: "RandConflictTest",
			input: input{
				filename:   "imports.go",
				interfaces: []string{"RandConflict"},
			},
			result: result{
				expectedImports: map[string]string{
					"math/rand": "rand2",
					"github.com/pm1381/sirish/internal/visitors/test_samples/rand": "rand",
				},
			},
		},
	}
	for _, each := range scenarios {
		t.Run(each.name, func(t *testing.T) {
			abs := internal.GetTestPathHelper(each.input.filename, "visitors")
			dt := new(dto.Types)
			for _, intr := range each.input.interfaces {
				err := dt.Set(intr)
				require.NoError(t, err)
			}
			wrapper := NewTypeVisitor(abs, internal.GenerateUniqueValues(*dt, nil))
			err := wrapper.Traverse()
			require.NoError(t, err)

			// assert
			assert.True(t, reflect.DeepEqual(each.result.expectedImports, wrapper.GetImports()))
		})
	}
}
