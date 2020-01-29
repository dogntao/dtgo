package utils

import (
	"reflect"
)

// sturct 转换为 map
func StructToMap(s interface{}) (retMap map[string]interface{}) {
	rt := reflect.TypeOf(s).Elem()
	rv := reflect.ValueOf(s).Elem()
	retMap = make(map[string]interface{}, 0)
	for i := 0; i < rt.NumField(); i++ {
		retMap[rt.Field(i).Tag.Get("json")] = rv.Field(i).Interface()
	}
	return
}

// struct arr 转换为 map arr
func StructArrToMapArr(s interface{}) (retMaps []map[string]interface{}) {
	sVal := reflect.ValueOf(s)
	sValKind := sVal.Kind().String()
	var retMap map[string]interface{}
	if sValKind == "slice" {
		for i := 0; i < sVal.Len(); i++ {
			if sVal.Index(i).Kind().String() == "struct" {
				retMap = make(map[string]interface{}, 0)
				sType := sVal.Index(i).Type()
				for v := 0; v < sVal.Index(i).NumField(); v++ {
					retMap[sType.Field(v).Tag.Get("json")] = sVal.Index(i).Field(v).Interface()
				}
				retMaps = append(retMaps, retMap)
			}
		}
	}
	return
}
