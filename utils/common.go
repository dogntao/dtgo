package utils

import (
	"crypto/md5"
	"encoding/hex"
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

// 过滤slice重复值
func ArrUniq(arrs []string) (res []string) {
	if len(arrs) > 0 {
		arrMaps := make(map[string]string)
		for _, v := range arrs {
			if v != "" {
				if _, ok := arrMaps[v]; !ok {
					arrMaps[v] = v
					res = append(res, v)
				}
			}
		}
	}
	return res
}

// arr 转换为 map
func ArrToMap(s []map[string]string, field string) (retMap map[string]interface{}) {
	if len(s) > 0 {
		retMap = make(map[string]interface{})
		for _, v := range s {
			if vVal, vOk := v[field]; vOk {
				if _, ok := retMap[vVal]; !ok {
					retMap[vVal] = v
				}
			}
		}
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

// Md5Fun
func Md5Fun(str string) (md5str string) {
	h := md5.New()
	h.Write([]byte(str))
	md5str = hex.EncodeToString(h.Sum(nil))
	return
}
