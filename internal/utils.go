package internal

import (
	"github.com/pm1381/sirish/internal/dto"
	"log"
	"os"
	"path/filepath"
)

func GetTestPathHelper(file string, rPath string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	abs := filepath.Join(wd, "test_samples", file)
	return abs
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
