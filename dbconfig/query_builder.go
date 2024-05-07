package dbconfig

import (
	"fmt"
	"reflect"
	"strings"
)

func QueryBuilder(payload interface{}, modelName string) string {
	reflectVal := reflect.ValueOf(payload)
	reflectType := reflect.TypeOf(payload)

	query := "INSERT INTO " + modelName

    onDuplicateFields := []string{}
	for i := 0; i < reflectType.NumField(); i++ {
        tempField :=  reflectType.Field(i).Tag.Get("cql")
        tempVal := fmt.Sprintf("%v", reflectVal.Field(i).Interface())
        if tempField == "created_at" || tempField == "updated_at" || tempField == "" || tempVal == "" {
            continue
        }

		if reflectVal.Field(i).Kind() == reflect.String {
            onDuplicateFields = append(onDuplicateFields, tempField+" = '"+tempVal+"'")
			continue
		}
        onDuplicateFields = append(onDuplicateFields, tempField+" = "+tempVal)
	}

	query += " SET " + strings.Join(onDuplicateFields, ",") + " ON DUPLICATE KEY UPDATE " + strings.Join(onDuplicateFields, ",") + " ;"
    return query
}

func QueryChecker(payload interface{}, modelName string) string {
	if payload == nil || modelName == "" {
		return ""
	}
	query := "SELECT ID FROM " + modelName + " WHERE id IN "
	reflectVal := reflect.ValueOf(payload).Interface().([]int32)

	listId := []string{}
	for i := 0; i < len(reflectVal); i++ {
		listId = append(listId, fmt.Sprintf("%v", reflectVal[i]))
	}
	query += "(" + strings.Join(listId, ",") + ");"
	return query
}
