package util

import (
	"github.com/fatih/structs"
	"strconv"
)

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

func ConvertDBKey(dbKey interface{}) interface{} {

	switch dbKey.(type) {
	case string:
		if _, err := strconv.Atoi(dbKey.(string)); err == nil {
			s, _ := strconv.ParseUint(dbKey.(string), 10, 32)
			return uint32(s)

		}
	}

	return dbKey
}
