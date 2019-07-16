package config

import (
	"reflect"
	"strconv"
)

func mergeCommandLineWithConfigFile(cl *Config, cf map[string]interface{}) error {
	t := reflect.TypeOf(cl)
	v := reflect.ValueOf(cl)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < t.NumField(); i++ {
		st := t.Field(i)
		plutoConfig := st.Tag.Get("pluto_config")
		m, ok := cf[plutoConfig]
		if !ok || reflect.TypeOf(m).Kind() != reflect.Map {
			continue
		}
		pc := v.FieldByName(st.Name).Interface()
		err := mergePlutoConfigWithMap(pc, m.(map[string]interface{}))
		return err
	}
	return nil
}

func mergePlutoConfigWithMap(pc interface{}, m map[string]interface{}) error {
	t := reflect.TypeOf(pc)
	v := reflect.ValueOf(pc)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		st := t.Field(i)
		plutoConfig := st.Tag.Get("pluto_config")
		if plutoConfig != "" {
			pc := v.FieldByName(st.Name).Interface()
			if mm, ok := m[plutoConfig]; !ok || reflect.TypeOf(mm).Kind() != reflect.Map {
				continue
			}
			err := mergePlutoConfigWithMap(pc, m[plutoConfig].(map[string]interface{}))
			return err
		}

		plutoValue := st.Tag.Get("pluto_value")
		if plutoValue != "" {
			i := v.FieldByName(st.Name).Interface()
			pv := i.(PlutoValue)
			value, ok := m[plutoValue]
			if !ok {
				continue
			}
			s := ""
			switch reflect.TypeOf(value).Kind() {
			case reflect.String:
				s = value.(string)
			case reflect.Int:
				s = strconv.Itoa(value.(int))
			case reflect.Int64:
				s = strconv.Itoa(int(value.(int64)))
			}
			err := pv.Set(s)
			return err
		}
	}
	return nil
}
