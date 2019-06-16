package objectTo

import (
	"reflect"
)

//struc转为map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
/*
---------------------
作者：heart_java
来源：CSDN
原文：https://blog.csdn.net/qq_29447481/article/details/72874847
版权声明：本文为博主原创文章，转载请附上博文链接！
*/

func B2String(bs []int8) string {
　　ba := make([]byte,0)
　　for _, b := range bs {
　　　　ba = append(ba, byte(b))
　　}
　　return string(ba)
}