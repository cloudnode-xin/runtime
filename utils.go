package runtime

import (
	"fmt"
	"reflect"
	"strings"
)

func getName(s Servicer) string {
	n := s.Name()
	if len(n) > 0 {
		return n
	}

	t := reflect.TypeOf(s)
	return fmt.Sprintf("(%s)", t.Elem().Name())
}

func getPath(s *Scope) string {
	scope, ok := s.parent.(*Scope)
	if !ok {
		return s.name
	}

	return fmt.Sprintf("%s.%s", getPath(scope), s.name)
}

func isReserve(name string) bool {
	return strings.HasPrefix(name, "(")
}
