package dto

import "strings"

type Types []string

func (t *Types) String() string {
	return "[" + strings.Join(*t, ",") + "]"
}

func (t *Types) Set(s string) error {
	if strings.Contains(s, ",") {
		elements := strings.Split(s, ",")
		*t = append(*t, elements...)
		return nil
	}
	*t = append(*t, s)
	return nil
}

func (t *Types) Exists(name string) bool {
	for _, eachType := range *t {
		if eachType == name {
			return true
		}
	}
	return false
}
