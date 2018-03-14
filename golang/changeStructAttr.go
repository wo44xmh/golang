package main

import (
	"reflect"
)

type MgzStrMap map[string]interface{}
type MgzBaseIntData map[string]map[int]interface{}
type Data struct {
	Id int `json:"-"`
	Title string `json:"title"`
	IdStr string `json:"id"`
}

func main()  {
	columns := make(MgzStrMap)
	columns["id"] = MgzStrMap{
		"title":"id",
		"output": get_map_partner,
		"attr":"Id",
	}
	partnerData := map[int]interface{}{1:"test"}

	Format_data([]Data{1,"test"}, columns, MgzBaseIntData{"partner" : partnerData})
}

func  Format_data(dataParam interface{}, attributesParam interface{}, params interface{}) {
	if attributes, ok := attributesParam.(map[string]interface{}); ok {
		dataInfo := reflect.ValueOf(dataParam)
		//dataType := dataInfo.Index(0).Type()
		if dataInfo.Kind() == reflect.Slice {
			l := reflect.Indirect(dataInfo).Len()
			//data := reflect.MakeSlice(dataInfo.Type(), l, l)
			for i := 0; i < l; i++ {
				val := reflect.Indirect(dataInfo).Index(i)

				for attrKey, attrVal := range attributes {
					attrValMap := attrVal.(MgzStrMap)
					if attrValMap["alias"] != nil {
						attrKey = attrValMap["alias"].(string)
					}

					originKey := attrKey
					if attrValMap["attr"] != nil {
						attrKey = attrValMap["attr"].(string)
					}

					if output, ok := attrValMap["output"]; ok {
						var has bool = false
						var temp reflect.Value
						var refVal reflect.Value
						if refVal = val.FieldByName(attrKey); refVal.IsValid() {
							has = true
							temp = reflect.Indirect(refVal)
						}

						outputTypeStr := reflect.TypeOf(output).Kind()
						if temp.IsValid() && outputTypeStr == reflect.Func && has {
							outputFunc := output.(func(interface{}, interface{}) interface{})
							var funcRes interface{}
							if temp.Kind() == reflect.Int {
								funcRes = outputFunc(int(temp.Int()), params)

							} else if temp.Kind() == reflect.Uint {
								funcRes = outputFunc(int(temp.Uint()), params)

							} else if temp.Kind() == reflect.String {
								funcRes = outputFunc(temp.String(), params)
							}
							SetModelAttribute(val, originKey, &funcRes)
						}
					}

				}
			}
		}

	}

}

func SetModelAttribute(model reflect.Value, key string, values interface{}) {
	l := model.NumField()
	//t := reflect.TypeOf(model).Elem()
	val := reflect.ValueOf(values)
	for i := 0; i < l; i++ {
		tagName := model.Type().Field(i).Tag.Get("json")
		if tagName != key {
			continue
		}
		if val.IsValid() {
			if val.Kind() == reflect.Ptr {
				if !val.IsNil() {
					//fmt.Println(reflect.Indirect(val).Elem().String())
					if reflect.Indirect(val).Elem().Kind() == reflect.Int {
						model.FieldByName(model.Type().Field(i).Name).SetInt(reflect.Indirect(val).Elem().Int())
					} else if reflect.Indirect(val).Elem().Kind() == reflect.String {
						model.FieldByName(model.Type().Field(i).Name).SetString(reflect.Indirect(val).Elem().String())
					}
				}
			}
		}
	}
}

func get_map_partner(valueParam interface{}, originParams interface{}) interface{} {
	if value, ok := valueParam.(int); ok {
		if params, ok := originParams.(MgzBaseIntData); ok {
			if temp, ok := params["partner"]; ok {
				if temp[value] != nil {
					resStr := temp[value].(string)
					return resStr;
				}
			}
		}
	}
	return "未知"
}