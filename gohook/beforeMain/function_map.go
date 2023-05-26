package beforeMain

import "reflect"

var Functions_Map map[string]interface{}

func BuildMap(functions []string) {
	for _, funcName := range functions {
		funcPtr := reflect.ValueOf(funcName).Pointer()
		Functions_Map[funcName] = funcPtr
	}
}
