package visitors

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/pm1381/sirish/internal/dto"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"math/big"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const smallLetters = "abcdefghijklmnopqrstuvwxyz"

func ExprToString(fset *token.FileSet, expr ast.Expr) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, expr)
	if err != nil {
		return "" // fallback: something is wrong, but avoid panicking in a generator
	}
	return buf.String()
}

func GenerateUniqueValues(primaryTargets dto.Types, secondaryTargets dto.Types) dto.Types {
	var interfacesMap = make(map[string]struct{})
	res := new(dto.Types)
	for _, target := range primaryTargets {
		*res = append(*res, target)
		interfacesMap[target] = struct{}{}
	}
	for _, target := range secondaryTargets {
		if _, ok := interfacesMap[target]; !ok {
			*res = append(*res, target)
			interfacesMap[target] = struct{}{}
		}
	}
	return *res
}

func MethodParamSnowflake(methodName string, paramIndex int, nameIndex int, length int, exactKey string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[cryptoRandSecure(int64(len(letters)))]
	}
	return fmt.Sprintf("%s%s%s_%d_%d",
		methodName,
		exactKey, // for example, it can be Ctx
		string(b),
		paramIndex,
		nameIndex,
	)
}

func cryptoRandSecure(max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		log.Println(err)
	}
	return nBig.Int64()
}
