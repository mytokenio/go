package mysql

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
)

func structMap(value reflect.Value) map[string]reflect.Value {
	m := make(map[string]reflect.Value)
	structValue(m, value)
	return m
}

func structValue(m map[string]reflect.Value, value reflect.Value) {
	if value.Type().Implements(reflect.TypeOf((*driver.Valuer)(nil)).Elem()) {
		return
	}
	switch value.Kind() {
	case reflect.Ptr:
		if value.IsNil() {
			return
		}
		structValue(m, value.Elem())
	case reflect.Struct:
		t := value.Type()
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).PkgPath != "" && !t.Field(i).Anonymous {
				continue
			}
			dbTag := t.Field(i).Tag.Get("db")
			switch dbTag {
			case dbTagDiscard:
				continue
			case dbTagEmpty:
				dbTag = t.Field(i).Name
			}

			if _, ok := m[dbTag]; !ok {
				m[dbTag] = value.Field(i)
			}
			structValue(m, value.Field(i))
		}
	}
}

func findPtr(column []string, value reflect.Value) ([]interface{}, error) {
	var dummy interface{}

	if value.Addr().Type().Implements(reflect.TypeOf((*sql.Scanner)(nil)).Elem()) {
		return []interface{}{value.Addr().Interface()}, nil
	}

	switch value.Kind() {
	case reflect.Struct:
		var ptr []interface{}
		m := structMap(value)
		for _, key := range column {
			if val, ok := m[key]; ok {
				ptr = append(ptr, val.Addr().Interface())
			} else {
				ptr = append(ptr, &dummy)
			}
		}
		return ptr, nil
	case reflect.Ptr:
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		return findPtr(column, value.Elem())
	}
	return []interface{}{value.Addr().Interface()}, nil
}

func struct2Map(data interface{}) (map[string]interface{}, error) {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	if t.Kind() != reflect.Ptr {
		return nil, errParamsBad
	}

	tElem := t.Elem()
	vElem := v.Elem()
	num := tElem.NumField()
	mapping := make(map[string]interface{})

	for i := 0; i < num; i++ {
		dbTag := tElem.Field(i).Tag.Get("db")
		if dbTag == dbTagDiscard {
			continue
		}
		if dbTag == dbTagEmpty {
			dbTag = tElem.Field(i).Name
		}

		value := vElem.Field(i).Interface()
		switch value.(type) {
		case NullString:
			value = vElem.Field(i).Interface().(NullString).String
		case sql.NullString:
			value = vElem.Field(i).Interface().(sql.NullString).String
		case NullBool:
			value = vElem.Field(i).Interface().(NullBool).Bool
		case sql.NullBool:
			value = vElem.Field(i).Interface().(sql.NullBool).Bool
		case NullInt64:
			value = vElem.Field(i).Interface().(NullInt64).Int64
		case sql.NullInt64:
			value = vElem.Field(i).Interface().(sql.NullInt64).Int64
		case NullFloat64:
			value = vElem.Field(i).Interface().(NullFloat64).Float64
		case sql.NullFloat64:
			value = vElem.Field(i).Interface().(sql.NullFloat64).Float64
		}

		mapping[dbTag] = value
	}

	return mapping, nil
}

func getColumns(data interface{}) ([]string, error) {
	t := reflect.TypeOf(data)

	switch t.Kind() {
	case reflect.Ptr:
		tElem := t.Elem()
		num := tElem.NumField()
		columns := make([]string, 0, num)
		for i := 0; i < num; i++ {
			dbTag := tElem.Field(i).Tag.Get("db")
			if dbTag == dbTagDiscard {
				continue
			} else if dbTag == dbTagEmpty {
				dbTag = tElem.Field(i).Name
			}
			columns = append(columns, dbTag)
		}
		return columns, nil
	case reflect.Map:
		switch data.(type) {
		case map[string]interface{}:
			columns := make([]string, 0, len(data.(map[string]interface{})))
			for column, _ := range data.(map[string]interface{}) {
				columns = append(columns, column)
			}
			return columns, nil
		}
	}

	return nil, errTypeInvalid
}

func getValuesFromReflect(v reflect.Value, t reflect.Type) []interface{} {
	tElem := t.Elem()
	vElem := v.Elem()
	num := tElem.NumField()
	values := make([]interface{}, 0, num)

	for i := 0; i < num; i++ {
		dbTag := tElem.Field(i).Tag.Get("db")
		value := vElem.Field(i).Interface()

		switch dbTag {
		case dbTagDiscard:
			continue
		case dbTagEmpty:
			values = append(values, value)
		default:
			values = append(values, value)
		}
	}

	return values
}

func getValues(data interface{}) ([]interface{}, error) {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	switch t.Kind() {
	case reflect.Ptr:
		tElem := t.Elem()
		vElem := v.Elem()
		num := tElem.NumField()
		values := make([]interface{}, 0, num)
		for i := 0; i < num; i++ {
			dbTag := tElem.Field(i).Tag.Get("db")
			if dbTag == dbTagDiscard {
				continue
			} else {
				values = append(values, vElem.Field(i).Interface())
			}
		}
		return values, nil
	case reflect.Map:
		switch data.(type) {
		case map[string]interface{}:
			values := make([]interface{}, 0, len(data.(map[string]interface{})))
			for _, value := range data.(map[string]interface{}) {
				values = append(values, value)
			}
			return values, nil
		}
	}

	return nil, errTypeInvalid
}
