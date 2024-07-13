package ruleEngine

import (
	"car-comparison-service/orm"
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

func GetCacheValueHelper[T any](re *RuleEngineExecutor, key string) (*T, error) {
	if !re.HasKey(key) {
		return nil, fmt.Errorf("value does not exist for %s", key)
	}

	if val, ok := re.GetValue(key).(T); !ok {
		return nil, fmt.Errorf("not a valid Type, key : %s", key)
	} else {
		return &val, nil
	}
}

func setPtr(val reflect.Value, value interface{}) error {
	if val.IsNil() {
		val.Set(reflect.New(val.Type().Elem()))
	}
	return setValue(val.Elem(), value)
}

func setValue(val reflect.Value, value interface{}) error {
	switch val.Kind() {
	case reflect.Ptr:
		return setPtr(val, value)
	}
	if !val.CanSet() {
		return fmt.Errorf("can't set value for struct field %s", val.Type())
	}
	if val.Type() != reflect.TypeOf(value) {
		return fmt.Errorf("type of field and value mismatch: %s - %s", val.Kind(), reflect.TypeOf(value).Kind())
	}
	val.Set(reflect.ValueOf(value))
	return nil
}

func SetDbForEngine(re *RuleEngineExecutor, db *gorm.DB) {
	re.SetValue("db", db)
}

func GetDbWithError(re *RuleEngineExecutor) (*gorm.DB, error) {
	val, err := GetCacheValueHelper[*gorm.DB](re, "db")
	if err != nil {
		return nil, err
	}
	return (*val).WithContext(re.Ctx), nil
}

func DbExecuteFunc(re *RuleEngineExecutor, result interface{}) error {
	qe, err := GetCacheValueHelper[orm.QueryEngine](re, "query")
	if err != nil {
		return err
	}
	return qe.Execute(result)
}
