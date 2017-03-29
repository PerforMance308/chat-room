package util

import "github.com/fatih/structs"

func MergeSimpleStruct(dst interface{}, src interface{}) {
	srcM := structs.Map(src)
	srcS := structs.New(src)
	dstS := structs.New(dst)

	for name, value := range srcM {
		if !srcS.Field(name).IsZero() {
			dstS.Field(name).Set(value)
		}
	}
}
